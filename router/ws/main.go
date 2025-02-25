package ws

import (
	"fmt"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/room"
	"github.com/victoroliveirab/settlers/router/ws/state"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func Handler(conn *types.WebSocketConnection, message *types.WebSocketMessage) error {
	handled, err := room.HandleMessage(conn, message)
	if handled {
		return err
	}
	return nil
}

func ClosePlayerConnection(player *types.GamePlayer) {
	room := state.RoomByID[player.Room]
	_, isPartOfGame := state.GameByRoom[player.Room]

	// Game hasn't started already, so just remove the player
	if !isPartOfGame {
		for i, entry := range room.Participants {
			if entry.Player.ID == player.ID {
				room.Participants[i] = types.RoomEntry{}
				break
			}
		}
		utils.BroadcastMessage(room, &types.WebSocketMessage{
			Type: "room.player-left",
			Payload: map[string]interface{}{
				"message": fmt.Sprintf("Player %s left the room", player.Username),
			},
		}, func(player *types.GamePlayer, err error) {
			logger.LogError(player.ID, "BroadcastMessage", -1, err)
			ClosePlayerConnection(player)
		})
		return
	}

	utils.BroadcastMessage(room, &types.WebSocketMessage{
		Type: "connection.lost",
		Payload: map[string]interface{}{
			"message": "connection lost",
		},
	}, func(player *types.GamePlayer, err error) {
		logger.LogError(player.ID, "BroadcastMessage", -1, err)
		ClosePlayerConnection(player)
	})
}
