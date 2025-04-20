package types

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RequestType string
type ResponseType string

type WebSocketConnection struct {
	Instance *websocket.Conn
	Mutex    sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWebSocketConnection(instance *websocket.Conn) *WebSocketConnection {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketConnection{
		Instance: instance,
		ctx:      ctx,
		cancel:   cancel,
	}
}

const (
	pingInterval = 10 * time.Second
	pongWait     = 15 * time.Second
	writeWait    = 5 * time.Second
)

func (c *WebSocketConnection) StartHeartBeat(onPingFail func()) {
	conn := c.Instance

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.Mutex.Lock()
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := conn.WriteMessage(websocket.PingMessage, nil)
				c.Mutex.Unlock()

				if err != nil {
					onPingFail()
					return
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *WebSocketConnection) Close() {
	c.cancel()
}

type WebSocketClientRequest struct {
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebSocketServerResponse struct {
	Type    ResponseType `json:"type"`
	Payload interface{}  `json:"payload"`
}
