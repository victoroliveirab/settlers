package entities

import (
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func NewPlayer(connection *types.WebSocketConnection, user *models.User, room *Room, onDisconnect func(player *GamePlayer)) *GamePlayer {
	return &GamePlayer{
		ID:           user.ID,
		Username:     user.Username,
		Connection:   connection,
		Color:        nil,
		Room:         room,
		OnDisconnect: onDisconnect,
	}
}

func (player *GamePlayer) ListenIncomingMessages(enqueueMessage func(msg *types.WebSocketClientRequest)) {
	defer player.OnDisconnect(player)
	for {
		parsedMessage, err := utils.ReadJson(player.Connection, player.ID)
		if err != nil {
			// DO PROPER HANDLING
			break
		}
		enqueueMessage(parsedMessage)
	}
}
