package entities

import (
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func NewPlayer(connection *types.WebSocketConnection, user *models.User) *GamePlayer {
	return &GamePlayer{
		ID:         user.ID,
		Username:   user.Username,
		Connection: connection,
		Color:      "",
		Room:       "",
	}
}

func (player *GamePlayer) ListenIncomingMessages(enqueueMessage func(msg *types.WebSocketMessage)) {
	defer player.Connection.Instance.Close()
	for {
		parsedMessage, err := utils.ReadJson(player.Connection, player.ID)
		if err != nil {
			// DO PROPER HANDLING
			break
		}
		enqueueMessage(parsedMessage)
	}
}
