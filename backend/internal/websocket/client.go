package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageProcessor interface {
	ProcessIncomingMessage(ctx context.Context, userID uuid.UUID, message []byte)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене здесь должна быть проверка домена
	},
}

func HandleWS(hub *Hub, processor MessageProcessor, w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	hub.Register <- client

	go writePump(client, hub)
	go readPump(client, hub, processor)
}

func readPump(client *Client, hub *Hub, processor MessageProcessor) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		processor.ProcessIncomingMessage(context.Background(), client.UserID, message)
	}
}

func writePump(client *Client, hub *Hub) {
	defer client.Conn.Close()

	for {
		message, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		client.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
