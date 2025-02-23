package manager

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/logger"
)

var color [4]string = [4]string{"green", "orange", "blue", "black"}

func CreateGameRoom(name, mapName string, capacity int) error {
	_, exists := maps.RoomByID[name]
	if exists {
		err := fmt.Errorf("Cannot create room #%s: already exists", name)
		return err
	}

	if mapName != "base4" {
		err := fmt.Errorf("Not supported map. Maps supported: base4")
		return err
	}

	players := make([]RoomEntry, capacity)

	lobby := Room{
		ID:       name,
		Capacity: capacity,
		MapName:  mapName,
		Players:  players,
	}
	maps.RoomByID[name] = lobby
	logger.LogMessage("system", "ws.manager.CreateGameRoom", "Created!")
	return nil
}

func AddPlayerToGameRoom(name string, userID int, conn *websocket.Conn) error {
	room, exists := maps.RoomByID[name]
	if exists {
		err := fmt.Errorf("Cannot join room #%s: room doesn't exist", name)
		return err
	}

	for i, player := range room.Players {
		if player.UserID == 0 {
			room.Players[i] = RoomEntry{
				Connection: conn,
				Color:      color[i],
				UserID:     userID,
				Bot:        false,
			}
			return nil
		}
	}

	err := fmt.Errorf("Cannot join room #%s: room full", name)
	return err
}
