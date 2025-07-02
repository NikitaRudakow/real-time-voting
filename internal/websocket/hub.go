package websocket

import (
	"sync"

	"real-time-voting/internal/dto"

	"github.com/rs/zerolog"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan dto.PollUpdateEvent
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	logger     zerolog.Logger
}

func NewHub(logger zerolog.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan dto.PollUpdateEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			h.logger.Info().Int("total_clients", len(h.clients)).Msg("Client registered")

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			h.logger.Info().Int("total_clients", len(h.clients)).Msg("Client unregistered")

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) Broadcast(event dto.PollUpdateEvent) {
	h.broadcast <- event
}

func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}
