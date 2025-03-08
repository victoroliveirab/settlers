package entities

import (
	"fmt"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
	wsUtils "github.com/victoroliveirab/settlers/router/ws/utils"
	"github.com/victoroliveirab/settlers/utils"
)

var color [4]string = [4]string{"green", "orange", "blue", "black"}

func NewRoom(id, mapName string, capacity int, onDestroy func(room *Room)) *Room {
	return &Room{
		ID:                id,
		Capacity:          capacity,
		MapName:           mapName,
		Participants:      make([]RoomEntry, capacity),
		Owner:             "",
		Status:            "prematch",
		incomingMsgQueue:  make(chan IncomingMessage, 32), // buffer incoming messages
		broadcastMsgQueue: make(chan BroadcastMessage),    // process msg immediatly, one by one
		handlers:          make([]RoomIncomingMessageHandler, 0),
		onDestroy:         onDestroy,
		Game:              nil,
		Private:           true,
	}
}

func (room *Room) AddPlayer(player *GamePlayer) error {
	room.Lock()
	defer room.Unlock()
	for _, spot := range room.Participants {
		if spot.Player != nil {
			if spot.Player.ID == player.ID {
				err := fmt.Errorf("Cannot add player %s to room#%s: already in room", player.Username, room.ID)
				return err
			}
		}
	}
	for i, spot := range room.Participants {
		if spot.Player == nil {
			player.Color = color[i]
			room.Participants[i] = RoomEntry{
				Player: player,
				Ready:  false,
				Bot:    false,
			}
			if room.Owner == "" {
				room.Owner = player.Username
			}
			return nil
		}
	}

	err := fmt.Errorf("Cannot join room #%s: room full", room.ID)
	return err
}

func (room *Room) ReconnectPlayer(playerID int64, connection *types.WebSocketConnection, onDisconnect func(player *GamePlayer)) (*GamePlayer, error) {
	room.Lock()
	defer room.Unlock()
	for index, spot := range room.Participants {
		if spot.Player != nil && spot.Player.ID == playerID {
			room.Participants[index].Player.Connection = connection
			room.Participants[index].Bot = false
			room.Participants[index].Player.OnDisconnect = onDisconnect
			return room.Participants[index].Player, nil
		}
	}
	err := fmt.Errorf("Cannot reconnect to room #%s: not part of room", room.ID)
	return nil, err
}

func (room *Room) TogglePlayerReadyState(playerID int64, newState bool) error {
	room.Lock()
	defer room.Unlock()

	for index, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			room.Participants[index].Ready = newState
			return nil
		}
	}

	err := fmt.Errorf("Cannot toggle player#%d ready state to %v: not part of room %s", playerID, newState, room.ID)
	return err
}

func (room *Room) RemovePlayer(playerID int64) error {
	room.Lock()
	defer room.Unlock()

	for index, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			if room.Game == nil {
				room.Participants[index] = RoomEntry{}
				if room.Owner == participant.Player.Username {
					err := room.assignNewOwner()
					if err != nil {
						// TODO: instead of destroying right away, schedule the destroy (10 seconds in the future)
						// So if the leaving of the last player was a refresh in the page, they don't lose the room
						room.Destroy(err.Error())
						return nil
					}
				}
				// REFACTOR: not ideal to have this message defined here and at handlers/pre-match/broadcast.go
				room.EnqueueBroadcastMessage(&types.WebSocketMessage{
					Type:    "room.new-update",
					Payload: room.ToMapInterface(),
				}, []int64{}, nil)
			} else {
				room.Participants[index].Player.Connection.Instance = nil
				room.Participants[index].Bot = true
				room.EnqueueBroadcastMessage(&types.WebSocketMessage{
					Type: "game.player-left",
					Payload: map[string]interface{}{
						"player": room.Participants[index].Player.Username,
						"bot":    true,
					},
				}, []int64{}, nil)
			}
			return nil
		}
	}

	err := fmt.Errorf("Cannot remove player#%d: not part of room %s", playerID, room.ID)
	return err
}

func (room *Room) ProgressStatus() error {
	if room.Status == "prematch" {
		room.Status = "setup"
		return nil
	}
	if room.Status == "setup" {
		room.Status = "match"
		return nil
	}
	err := fmt.Errorf("Cannot proceed status %s", room.Status)
	return err
}

func (room *Room) assignNewOwner() error {
	for _, participant := range room.Participants {
		if participant.Player != nil {
			room.Owner = participant.Player.Username
			return nil
		}
	}
	err := fmt.Errorf("Cannot assign a new owner to room %s: no players left", room.ID)
	return err
}

func (room *Room) RegisterIncomingMessageHandler(f RoomIncomingMessageHandler) {
	room.handlers = append(room.handlers, f)
}

func (room *Room) EnqueueIncomingMessage(player *GamePlayer, msg *types.WebSocketMessage) {
	room.incomingMsgQueue <- IncomingMessage{
		Room:    room,
		Player:  player,
		Message: msg,
	}
}

func (room *Room) ProcessIncomingMessages() {
	for {
		item := <-room.incomingMsgQueue
		message := item.Message
		sender := item.Player

		var handled bool
		var err error

		for _, handler := range room.handlers {
			handled, err = handler(sender, message)
			if handled || err != nil {
				break
			}
		}
		if handled && err == nil {
			continue
		}
		// TODO: handle error for appropriate player
	}
}

func (room *Room) EnqueueBroadcastMessage(msg *types.WebSocketMessage, excludedUserIDs []int64, onSend func()) {
	room.broadcastMsgQueue <- BroadcastMessage{
		ExcludedIDs: excludedUserIDs,
		Message:     msg,
		OnSend:      onSend,
	}
}

func (room *Room) ProcessBroadcastRequests() {
	for {
		item := <-room.broadcastMsgQueue
		msg := item.Message
		excludedUserIDs := item.ExcludedIDs
		onSendCb := item.OnSend

		room.Lock()
		for _, participant := range room.Participants {
			player := participant.Player
			if player == nil || player.Connection.Instance == nil || utils.SliceContains(excludedUserIDs, player.ID) {
				continue
			}

			err := wsUtils.WriteJson(player.Connection, player.ID, msg)

			if err != nil {
				fmt.Println("error for player ", player.ID, err)
				// TODO: handle error properly
			}
		}
		room.Unlock()
		if onSendCb != nil {
			go onSendCb()
		}
	}
}

func (room *Room) Destroy(reason string) {
	logger.LogSystemMessage("room.Destroy", reason)
	room.onDestroy(room)
}

func (room *Room) ToMapInterface() map[string]interface{} {
	return map[string]interface{}{
		"id":           room.ID,
		"capacity":     room.Capacity,
		"map":          room.MapName,
		"participants": room.Participants,
		"owner":        room.Owner,
		"status":       room.Status,
	}
}
