package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func ReconnectPlayer(room *entities.Room, playerID int64, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	requestType := types.RequestType("room.connect")
	player, err := room.ReconnectPlayer(playerID, conn, onDisconnect)
	if err != nil {
		logger.LogError(playerID, "prematch.ReconnectPlayer", -1, err)
		return nil, err
	}

	msg := BuildRoomMessage(room, fmt.Sprintf("%s.success", requestType))
	wsErr := utils.WriteJson(player.Connection, player.ID, msg)
	if wsErr != nil {
		return nil, wsErr
	}
	return player, nil
}
