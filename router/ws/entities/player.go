package entities

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/db/models"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func NewPlayer(user *models.User, room *Room, onDisconnect func(player *GamePlayer)) *GamePlayer {
	ctx, cancel := context.WithCancel(context.Background())
	return &GamePlayer{
		ID:           user.ID,
		Username:     user.Username,
		Color:        nil,
		Room:         room,
		OnDisconnect: onDisconnect,
		ctx:          ctx,
		cancelFunc:   cancel,
	}
}

func (player *GamePlayer) Connect(
	conn *websocket.Conn,
	enqueueMessage func(msg *types.WebSocketClientRequest),
) {
	player.connMu.Lock()
	defer player.connMu.Unlock()

	if player.activeConn != nil {
		logger.LogSystemMessage(
			fmt.Sprintf("player.%d.Connect()", player.ID),
			fmt.Sprintf("closing previous connection for player %s", player.Username),
		)
		player.activeConn.Close()
	}

	if player.connCancel != nil {
		logger.LogSystemMessage(
			fmt.Sprintf("player.%d.Connect()", player.ID),
			fmt.Sprintf("cancelling previous connection context for player %s", player.Username),
		)
		player.connCancel()
		player.connCancel = nil
	}

	if player.ctx.Err() != nil {
		logger.LogSystemMessage(
			fmt.Sprintf("player.%d.Connect()", player.ID),
			fmt.Sprintf("resetting  main context for player %s", player.Username),
		)
		player.ctx, player.cancelFunc = context.WithCancel(context.Background())
	}

	player.activeConn = conn

	connCtx, cancel := context.WithCancel(context.Background())
	player.connCancel = cancel

	go player.listenIncomingMessages(connCtx, enqueueMessage)
}

func (player *GamePlayer) GetConnection() *websocket.Conn {
	return player.activeConn
}

func (player *GamePlayer) Disconnect() {
	player.connMu.Lock()
	defer player.connMu.Unlock()

	logger.LogSystemMessage(
		fmt.Sprintf("player.%d.Disconnect()", player.ID),
		fmt.Sprintf("disconnecting player %s", player.Username),
	)

	if player.connCancel != nil {
		player.connCancel()
		player.connCancel = nil
	}

	if player.activeConn != nil {
		player.activeConn.Close()
		player.activeConn = nil
	}
}

func (player *GamePlayer) listenIncomingMessages(ctx context.Context, enqueueMessage func(msg *types.WebSocketClientRequest)) {
	logger.LogSystemMessage(
		fmt.Sprintf("listenIncomingMessages.%d", player.ID),
		fmt.Sprintf("starting listening for messages of player %s", player.Username),
	)

	defer func() {
		player.connMu.Lock()
		defer player.connMu.Unlock()

		logger.LogSystemMessage(
			fmt.Sprintf("listenIncomingMessages.%d.cleanup", player.ID),
			fmt.Sprintf("stopping listening for messages of player %s", player.Username),
		)

		if player.activeConn != nil {
			logger.LogSystemMessage(
				fmt.Sprintf("listenIncomingMessages.%d.cleanup", player.ID),
				fmt.Sprintf("cleaning up connection of player %s", player.Username),
			)
			player.activeConn = nil
			player.connCancel = nil
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.LogSystemMessage(
				fmt.Sprintf("listenIncomingMessages.%d.ctx.Done()", player.ID),
				fmt.Sprintf("connection context cancelled for player %s", player.Username),
			)
			return
		case <-player.ctx.Done():
			logger.LogSystemMessage(
				fmt.Sprintf("listenIncomingMessages.%d.player.ctx.Done()", player.ID),
				fmt.Sprintf("player context cancelled for player %s", player.Username),
			)
			return
		default:
		}

		m, message, err := player.activeConn.ReadMessage()

		select {
		case <-ctx.Done():
			return
		case <-player.ctx.Done():
			return
		default:
		}

		if err != nil {
			logger.LogError(player.ID, "conn.ReadMessage", m, err)
			return
		}

		var parsedMessage types.WebSocketClientRequest
		err = json.Unmarshal(message, &parsedMessage)
		if err != nil {
			logger.LogError(player.ID, "json.Unmarshal", -1, err)
			continue
		}

		enqueueMessage(&parsedMessage)
	}
}

func (player *GamePlayer) IsConnected() bool {
	player.connMu.Lock()
	defer player.connMu.Unlock()

	return player.activeConn != nil
}

func (player *GamePlayer) WriteJSON(message *types.WebSocketServerResponse) error {
	player.connMu.Lock()
	defer player.connMu.Unlock()

	if player.activeConn == nil {
		return fmt.Errorf("connection is nil")
	}

	err := player.activeConn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return fmt.Errorf("failed to set write deadline: %w", err)
	}

	return player.activeConn.WriteJSON(message)
}

func (player *GamePlayer) WriteJsonError(requestType types.RequestType, err error) error {
	message := &types.WebSocketServerResponse{
		Type: types.ResponseType(fmt.Sprintf("%s.error", requestType)),
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	}
	return player.WriteJSON(message)
}

func (player *GamePlayer) Destroy() {
	logger.LogSystemMessage(
		fmt.Sprintf("player.%d.Destroy()", player.ID),
		fmt.Sprintf("destroying player %s", player.Username),
	)
	player.Room = nil
}

func (player *GamePlayer) SetLastActiveTimestamp(timestamp time.Time) {
	player.LastTimeActive = timestamp
}
