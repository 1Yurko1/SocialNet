package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/websocket"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
	hub      *websocket.Hub
}

func NewChatService(cr *repository.ChatRepository, hub *websocket.Hub) *ChatService {
	return &ChatService{chatRepo: cr, hub: hub}
}

func (s *ChatService) GetOrCreatePrivateChat(ctx context.Context, user1, user2 uuid.UUID) (uuid.UUID, error) {
	chatID, err := s.chatRepo.GetOrCreatePrivateChat(ctx, user1, user2)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}

func (s *ChatService) ProcessIncomingMessage(ctx context.Context, senderID uuid.UUID, rawMsg []byte) {
	var wsMsg models.WSMessage
	if err := json.Unmarshal(rawMsg, &wsMsg); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return
	}

	if wsMsg.Type == "private_msg" {
		targetUserID, err := uuid.Parse(wsMsg.ToUserID)
		if err != nil {
			log.Printf("Invalid ToUserID received: %s. Error: %v", wsMsg.ToUserID, err)
			return
		}

		chatID, err := s.chatRepo.GetOrCreatePrivateChat(ctx, senderID, targetUserID)
		if err != nil {
			log.Printf("Error getting/creating chat: %v", err)
			return
		}

		msg := &models.Message{
			ID:        uuid.New(),
			ChatID:    chatID,
			SenderID:  senderID,
			Content:   wsMsg.Text,
			CreatedAt: time.Now(),
		}

		if err := s.chatRepo.SaveMessage(ctx, msg); err != nil {
			log.Printf("Error saving message to DB: %v", err)
		}

		// При ответе конвертируем UUID обратно в строку для JSON
		response, _ := json.Marshal(models.WSMessage{
			Type:      "private_msg",
			ChatID:    chatID.String(), // .String()
			Text:      wsMsg.Text,
			SenderID:  senderID.String(), // .String()
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		})

		s.hub.SendToUser(targetUserID, response)
	}
}

func (s *ChatService) GetChatHistory(ctx context.Context, chatID uuid.UUID) ([]*models.Message, error) {
	return s.chatRepo.GetMessagesByChat(ctx, chatID, 50, 0)
}
