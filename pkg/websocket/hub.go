package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *zerolog.Logger
}

func NewHub(logger *zerolog.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.logger.Info().Str("client", client.conn.RemoteAddr().String()).Msg("Client registered")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.logger.Info().Str("client", client.conn.RemoteAddr().String()).Msg("Client unregistered")
			}
		case message := <-h.broadcast:
			h.logger.Info().Int("client_count", len(h.clients)).Msg("Broadcasting message to clients")
			for client := range h.clients {
				select {
				case client.send <- message:
					h.logger.Info().Str("client", client.conn.RemoteAddr().String()).Msg("Message sent to client")
				default:
					close(client.send)
					delete(h.clients, client)
					h.logger.Warn().Str("client", client.conn.RemoteAddr().String()).Msg("Failed to send message, client removed")
				}
			}
		}
	}
}

func (h *Hub) BroadcastUpdate(data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error marshaling message")
		return
	}
	h.broadcast <- message
	h.logger.Info().Msg("Update broadcasted to all clients")
}

func (h *Hub) SendToClient(userID int, msgType string, msg string) {
	for client := range h.clients {
		if client.userID == userID {
			message := map[string]string{"type": msgType, "message": msg}
			messageJSON, _ := json.Marshal(message)
			client.send <- messageJSON
			h.logger.Info().Str("client", client.conn.RemoteAddr().String()).Msg("Message sent to specific client")
		}
	}
}

func (h *Hub) SendToClientError(errorMsg string, w http.ResponseWriter) {
	message := map[string]string{"type": "error", "message": errorMsg}
	messageJSON, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(messageJSON)
	h.logger.Error().Msg("Error message sent to client")
}
