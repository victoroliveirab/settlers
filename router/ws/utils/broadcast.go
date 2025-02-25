package utils

import (
	"sync"

	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BroadcastMessage(room *types.Room, message *types.WebSocketMessage, onError func(player *types.GamePlayer, err error)) {
	errors := make([]error, len(room.Participants))
	hadError := false
	var errorsMutex sync.Mutex
	go func() {
		var wg sync.WaitGroup
		for index, entry := range room.Participants {
			player := entry.Player
			if player.Connection == nil {
				continue
			}
			wg.Add(1)
			go func(conn *types.WebSocketConnection, userID int64) {
				defer wg.Done()
				err := WriteJson(conn, userID, message)
				if err != nil {
					errorsMutex.Lock()
					hadError = true
					errors[index] = err
					errorsMutex.Unlock()
				}
			}(player.Connection, player.ID)
		}
		wg.Wait()
	}()

	if !hadError {
		return
	}

	for index, error := range errors {
		onError(room.Participants[index].Player, error)
	}
}
