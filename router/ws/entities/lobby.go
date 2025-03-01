package entities

import (
	"fmt"

	"github.com/victoroliveirab/settlers/logger"
)

func NewLobby() *Lobby {
	return &Lobby{
		roomByPlayer: make(map[int64]*Room),
		rooms:        make(map[string]*Room),
	}
}

// CreateRoom generates a new room and stores it in the Hub
func (lobby *Lobby) CreateRoom(id, mapName string, capacity int) (*Room, error) {
	lobby.Lock()
	defer lobby.Unlock()
	_, exists := lobby.rooms[id]
	if exists {
		err := fmt.Errorf("Cannot create room #%s: already exists", id)
		return nil, err
	}

	if mapName != "base4" {
		err := fmt.Errorf("Not supported map. Maps supported: base4")
		return nil, err
	}

	newRoom := NewRoom(id, mapName, capacity, func(room *Room) {
		lobby.DestroyRoom(room.ID)
	})
	lobby.rooms[id] = newRoom
	return newRoom, nil
}

// GetRoom retrieves a room by ID
func (lobby *Lobby) GetRoom(id string) (*Room, bool) {
	lobby.Lock()
	defer lobby.Unlock()
	room, exists := lobby.rooms[id]
	return room, exists
}

func (lobby *Lobby) DestroyRoom(roomID string) error {
	room, exists := lobby.rooms[roomID]
	if !exists {
		err := fmt.Errorf("Cannot destroy room %s: no such room", roomID)
		return err
	}

	if room.Game != nil {
		err := fmt.Errorf("Cannot destroy room %s: game ongoing", roomID)
		return err
	}

	delete(lobby.rooms, roomID)
	logger.LogSystemMessage("lobby.DestroyRoom", fmt.Sprintf("Destroyed room %s", roomID))
	return nil
}
