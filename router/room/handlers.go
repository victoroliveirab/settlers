package room

import (
	"fmt"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/state"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func CreateGameRoom(name, mapName string, capacity int) error {
	_, exists := state.RoomByID[name]
	if exists {
		err := fmt.Errorf("Cannot create room #%s: already exists", name)
		return err
	}

	if mapName != "base4" {
		err := fmt.Errorf("Not supported map. Maps supported: base4")
		return err
	}

	players := make([]types.RoomEntry, capacity)

	lobby := &types.Room{
		ID:           name,
		Capacity:     capacity,
		MapName:      mapName,
		Participants: players,
	}
	state.RoomByID[name] = lobby
	logger.LogSystemMessage("ws.manager.CreateGameRoom", "Created!")
	return nil
}
