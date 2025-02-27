package entities

import "fmt"

func NewLobby() *Lobby {
	return &Lobby{
		rooms: make(map[string]*Room),
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

	newRoom := NewRoom(id, mapName, capacity)
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
