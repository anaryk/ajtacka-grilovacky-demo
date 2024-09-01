package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID int
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, logger *zerolog.Logger) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	hub.register <- client

	logger.Info().Msg("New WebSocket connection established")

	go client.readPump(logger)
	go client.writePump(logger)
}

func (c *Client) readPump(logger *zerolog.Logger) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		logger.Info().Msg("WebSocket connection closed")
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error().Err(err).Msg("WebSocket error")
			}
			break
		}
		logger.Info().Str("msg", string(message)).Msg("Received message from WebSocket client")
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump(logger *zerolog.Logger) {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				logger.Info().Msg("WebSocket send channel closed")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to write message to WebSocket client")
				return
			}
			logger.Info().Msg("Message successfully sent to WebSocket client")
		}
	}
}
