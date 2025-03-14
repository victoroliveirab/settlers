package entities

import (
	"fmt"
	"sort"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
	wsUtils "github.com/victoroliveirab/settlers/router/ws/utils"
	"github.com/victoroliveirab/settlers/utils"
)

var availableColors []string = []string{"palegreen", "orange", "maroon", "lemonchiffon", "blue", "crimson", "orangered", "white", "aliceblue", "lightslategray"}

func NewRoom(id, mapName string, capacity int, params RoomParams, onDestroy func(room *Room)) *Room {
	return &Room{
		ID:               id,
		Capacity:         capacity,
		MapName:          mapName,
		Participants:     make([]RoomEntry, capacity),
		Owner:            "",
		Colors:           availableColors,
		params:           params,
		Status:           "prematch",
		incomingMsgQueue: make(chan IncomingMessage, 32), // buffer incoming messages
		outgoingMsgQueue: make(chan OutgoingMessage),     // process msg immediatly, one by one
		handlers:         make([]RoomIncomingMessageHandler, 0),
		onDestroy:        onDestroy,
		Game:             nil,
		Private:          true,
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
			player.Color = availableColors[i]
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

func (room *Room) ChangePlayerColor(playerID int64, color string) error {
	room.Lock()
	defer room.Unlock()

	if !utils.SliceContains(availableColors, color) {
		err := fmt.Errorf("Cannot use color %s: unknown color", color)
		return err
	}

	for _, participant := range room.Participants {
		if participant.Player != nil && participant.Player.Color == color {
			err := fmt.Errorf("Cannot use color %s: color taken", color)
			return err
		}
	}
	for index, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			room.Participants[index].Player.Color = color
			return nil
		}
	}

	err := fmt.Errorf("Cannot set color for %d: player not found", playerID)
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
				// room.EnqueueBroadcastMessage(&types.WebSocketMessage{
				// 	Type:    "room.new-update",
				// 	Payload: room.ToMapInterface(),
				// }, []int64{}, nil)
			} else {
				room.Participants[index].Player.Connection.Instance = nil
				room.Participants[index].Bot = true
				// room.EnqueueBroadcastMessage(&types.WebSocketMessage{
				// 	Type: "game.player-left",
				// 	Payload: map[string]interface{}{
				// 		"player": room.Participants[index].Player.Username,
				// 		"bot":    true,
				// 	},
				// }, []int64{}, nil)
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

func (room *Room) EnqueueIncomingMessage(player *GamePlayer, msg *types.WebSocketClientRequest) {
	room.incomingMsgQueue <- IncomingMessage{
		Room:    room,
		Player:  player,
		Message: msg,
	}
}

func (room *Room) EnqueueOutgoingMessage(msg *types.WebSocketServerResponse, recipients []string, onSend func()) {
	room.outgoingMsgQueue <- OutgoingMessage{
		Callback:   onSend,
		Message:    msg,
		Recipients: recipients,
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

func (room *Room) ProcessOutgoingMessages() {
	for {
		item := <-room.outgoingMsgQueue
		var recipients []RoomEntry

		// Copy so if a participant disconnects mid broadcast and the room.Participants array changes, we don't panic
		room.Lock()
		for _, participant := range room.Participants {
			if len(item.Recipients) == 0 {
				recipients = append(recipients, participant)
				continue
			}
			for _, recipient := range item.Recipients {
				if participant.Player.Username == recipient {
					recipients = append(recipients, participant)
				}
			}
		}
		room.Unlock()

		for _, participant := range recipients {
			player := participant.Player
			if player == nil || player.Connection.Instance == nil {
				continue
			}
			wsErr := wsUtils.WriteJson(player.Connection, player.ID, item.Message)
			if wsErr != nil {
				fmt.Println("error for player ", player.ID, wsErr)
				// TODO: handle error here as well
				continue
			}
		}

		if item.Callback != nil {
			go item.Callback()
		}
	}
}

func (room *Room) Destroy(reason string) {
	logger.LogSystemMessage("room.Destroy", reason)
	room.onDestroy(room)
}

func (room *Room) Params() []RoomParamsMetaEntry {
	var entries []RoomParamsMetaEntry
	for _, v := range room.params.Meta {
		entries = append(entries, RoomParamsMetaEntry{
			Key:         v.Key,
			Description: v.Description,
			Priority:    v.Priority,
			Value:       room.params.Values[v.Key],
			Values:      v.Values,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Priority == entries[j].Priority {
			return entries[i].Key < entries[j].Key // Alphabetical order
		}
		return entries[i].Priority > entries[j].Priority // Higher priority first
	})

	return entries
}

func (room *Room) UpdateParam(player *GamePlayer, key string, value int) error {
	if room.Owner != player.Username {
		err := fmt.Errorf("cannot update param %s in room %s: not room owner", key, room.ID)
		return err
	}

	_, exists := room.params.Values[key]
	if !exists {
		err := fmt.Errorf("unknown param: %s", key)
		return err
	}
	if !utils.SliceContains(room.params.Meta[key].Values, value) {
		err := fmt.Errorf("invalid value for param %s: %d", key, value)
		return err
	}
	room.params.Values[key] = value
	return nil
}
