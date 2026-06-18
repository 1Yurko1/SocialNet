package websocket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client представляет одно активное соединение
type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
	Send   chan []byte // Канал для отправки сообщений клиенту
}

// Hub управляет всеми клиентами
type Hub struct {
	// Карта активных клиентов: UserID -> Client
	Clients map[uuid.UUID]*Client

	// Регистрация нового клиента
	Register chan *Client

	// Удаление клиента
	Unregister chan *Client

	// Канал для входящих сообщений (broadcast или private)
	Broadcast chan []byte

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			// В будущем здесь будет логика рассылки сообщения конкретным участникам чата
			h.mu.RLock()
			for _, client := range h.Clients {
				client.Send <- message
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser отправляет сообщение конкретному пользователю
func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.Clients[userID]; ok {
		client.Send <- message
	}
}
