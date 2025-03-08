package prematch

import (
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ConnectPlayer(room *entities.Room, user *models.User, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	player := entities.NewPlayer(conn, user, room, func(player *entities.GamePlayer) {
		room.RemovePlayer(user.ID)
	})
	err := room.AddPlayer(player)
	if err != nil {
		logger.LogError(player.ID, "prematch.ConnectPlayer", -1, err)
		return nil, err
	}

	err = SendConnectPlayerSuccess(player)
	if err != nil {
		return nil, err
	}
	return player, nil
}
