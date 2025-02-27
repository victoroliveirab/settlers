package prematch

import "github.com/victoroliveirab/settlers/router/ws/entities"

func roomToMapInterface(room *entities.Room) map[string]interface{} {
	return map[string]interface{}{
		"id":           room.ID,
		"capacity":     room.Capacity,
		"map":          room.MapName,
		"participants": room.Participants,
	}
}
