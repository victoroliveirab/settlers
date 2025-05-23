package entities

import (
	"fmt"

	mapsdefinitions "github.com/victoroliveirab/settlers/core/maps"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/utils"
)

func NewLobby() *Lobby {
	return &Lobby{
		availableRooms: []string{"base4"},
		rooms:          make(map[string]*Room),
	}
}

// CreateRoom generates a new room and stores it in the Hub
func (lobby *Lobby) CreateRoom(id, mapName string, capacity, randSeed int) (*Room, error) {
	lobby.Lock()
	defer lobby.Unlock()
	_, exists := lobby.rooms[id]
	if exists {
		err := fmt.Errorf("Cannot create room #%s: already exists", id)
		return nil, err
	}

	if !utils.SliceContains(lobby.availableRooms, mapName) {
		err := fmt.Errorf("Not supported map. Maps supported: %v", lobby.availableRooms)
		return nil, err
	}

	meta, err := mapsdefinitions.GetMetadata(mapName)
	if err != nil {
		return nil, err
	}

	paramsMeta := make(RoomParamsMeta)
	paramsValues := make(map[string]int)
	for key := range meta.Params {
		paramsMeta[key] = RoomParamsMetaEntry{
			Key:         key,
			Description: meta.Params[key].Description,
			Label:       meta.Params[key].Label,
			Priority:    meta.Params[key].Priority,
			Value:       meta.Params[key].Default,
			Values:      meta.Params[key].Values,
		}
		paramsValues[key] = meta.Params[key].Default
	}
	params := RoomParams{
		Meta:       paramsMeta,
		MaxPlayers: meta.Players.Max,
		MinPlayers: meta.Players.Min,
		Values:     paramsValues,
	}

	newRoom := NewRoom(id, mapName, capacity, randSeed, params, func(room *Room) {
		lobby.RemoveRoom(room.ID)
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

func (lobby *Lobby) RemoveRoom(roomID string) error {
	_, exists := lobby.rooms[roomID]
	if !exists {
		err := fmt.Errorf("Cannot remove room %s from lobby: no such room", roomID)
		return err
	}

	delete(lobby.rooms, roomID)
	logger.LogSystemMessage("lobby.RemoveRoom", fmt.Sprintf("Removed room %s from lobby", roomID))
	return nil
}
