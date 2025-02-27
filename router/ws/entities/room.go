package entities

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

var color [4]string = [4]string{"green", "orange", "blue", "black"}

func NewRoom(id, mapName string, capacity int) *Room {
	return &Room{
		ID:                id,
		Capacity:          capacity,
		MapName:           mapName,
		Participants:      make([]RoomEntry, capacity),
		incomingMsgQueue:  make(chan IncomingMessage, 32), // buffer incoming messages
		broadcastMsgQueue: make(chan BroadcastMessage),    // process msg immediatly, one by one
		handlers:          make([]RoomIncomingMessageHandler, 0),
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
				return nil
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
			return nil
		}
	}

	err := fmt.Errorf("Cannot join room #%s: room full", room.ID)
	return err
}

func (room *Room) TogglePlayerReadyState(playerID int64, newState bool) error {
	room.Lock()
	defer room.Unlock()

	for _, participant := range room.Participants {
		if participant.Player != nil && participant.Player.ID == playerID {
			participant.Ready = newState
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
			room.Participants[index] = RoomEntry{}
			return nil
		}
	}

	err := fmt.Errorf("Cannot remove player#%d: not part of room %s", playerID, room.ID)
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
			handled, err = handler(room, sender, message)
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

func (room *Room) EnqueueBroadcastMessage(msg *types.WebSocketMessage, excludedUserIDs []int64) {
	room.broadcastMsgQueue <- BroadcastMessage{
		ExcludedIDs: excludedUserIDs,
		Message:     msg,
	}
}

func (room *Room) ProcessBroadcastRequests() {
	for {
		item := <-room.broadcastMsgQueue
		msg := item.Message
		// excludedUserIDs := item.ExcludedIDs

		room.Lock()
		for _, participant := range room.Participants {
			player := participant.Player
			if player == nil || player.Connection.Instance == nil {
				continue
			}

			err := utils.WriteJson(player.Connection, player.ID, msg)

			if err != nil {
				fmt.Println("error for player ", player.ID, err)
				// TODO: handle error properly
			}
		}
		room.Unlock()

	}
}
