package types

import (
	"encoding/json"
)

type RequestType string
type ResponseType string

// const (
// 	pingInterval = 10 * time.Second
// 	pongWait     = 15 * time.Second
// 	writeWait    = 5 * time.Second
// )

// func (c *WebSocketConnection) StartHeartBeat(onPingFail func()) {
// 	conn := c.Instance
//
// 	conn.SetReadDeadline(time.Now().Add(pongWait))
// 	conn.SetPongHandler(func(string) error {
// 		conn.SetReadDeadline(time.Now().Add(pongWait))
// 		return nil
// 	})
//
// 	ticker := time.NewTicker(pingInterval)
//
// 	go func() {
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ticker.C:
// 				c.Mutex.Lock()
// 				conn.SetWriteDeadline(time.Now().Add(writeWait))
// 				err := conn.WriteMessage(websocket.PingMessage, nil)
// 				c.Mutex.Unlock()
//
// 				if err != nil {
// 					onPingFail()
// 					return
// 				}
//
// 				conn.SetWriteDeadline(time.Time{})
// 			case <-c.ctx.Done():
// 				return
// 			}
// 		}
// 	}()
// }
//
// func (c *WebSocketConnection) Close() {
// 	c.cancel()
// }

type WebSocketClientRequest struct {
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebSocketServerResponse struct {
	Type    ResponseType `json:"type"`
	Payload interface{}  `json:"payload"`
}
