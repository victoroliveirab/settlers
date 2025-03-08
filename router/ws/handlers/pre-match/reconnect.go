package prematch

import (
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	player, err := room.ReconnectPlayer(playerID, conn, onDisconnect)
	if err != nil {
		logger.LogError(playerID, "prematch.ReconnectPlayer", -1, err)
		return nil, err
	}

	err = SendConnectPlayerSuccess(player)
	if err != nil {
		return nil, err
	}
	return player, nil
}
