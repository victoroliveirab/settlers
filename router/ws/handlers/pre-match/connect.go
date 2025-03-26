package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func ConnectPlayer(room *entities.Room, user *models.User, conn *types.WebSocketConnection, onDisconnect func(player *entities.GamePlayer)) (*entities.GamePlayer, error) {
	requestType := types.RequestType("room.connect")
	player := entities.NewPlayer(conn, user, room, func(player *entities.GamePlayer) {
		room.RemovePlayer(user.ID)
	})
	err := room.AddPlayer(player)
	if err != nil {
		logger.LogError(player.ID, "prematch.ConnectPlayer", -1, err)
		return nil, err
	}

	msg := BuildRoomMessage(room, fmt.Sprintf("%s.success", requestType))
	wsErr := utils.WriteJson(player.Connection, player.ID, msg)
	if wsErr != nil {
		return nil, wsErr
	}
	return player, nil
}
