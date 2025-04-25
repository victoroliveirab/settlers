package postmatch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection) (*entities.GamePlayer, error) {
	player, err := room.ReconnectPlayer(playerID, conn, func(player *entities.GamePlayer) {
		room.RemovePlayer(playerID)
	})
	if err != nil {
		wsErr := utils.WriteJsonError(conn, playerID, "match.ReconnectPlayer", err)
		return nil, wsErr
	}
	game := room.Game
	if game == nil {
		err := fmt.Errorf("Game not assigned to the room", room.ID)
		return nil, err
	}

	msg := BuildPostMatchHydrateMessage(room)
	wsErr := utils.WriteJson(player.Connection, player.ID, msg)
	return player, wsErr
}
