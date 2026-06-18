package handlers

import (
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: s}
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	// 1. Извлекаем ID чата из параметров URL (:id)
	chatIDStr := c.Param("id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id format",
		})
		return
	}

	// 2. Вызываем бизнес-логику сервиса
	// Мы передаем контекст запроса, чтобы сервис мог отменить операцию, если клиент отключится
	messages, err := h.chatService.GetChatHistory(c.Request.Context(), chatID)
	if err != nil {
		// Логируем ошибку для разработчика, а пользователю отдаем общую ошибку
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve chat history",
		})
		return
	}

	// 3. Возвращаем список сообщений
	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) CreatePrivateChat(c *gin.Context) {
	// Получаем ID текущего пользователя из JWT (установлен в AuthMiddleware)
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req struct {
		TargetUserID string `json:"target_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_user_id is required"})
		return
	}

	targetID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target user id"})
		return
	}

	chatID, err := h.chatService.GetOrCreatePrivateChat(c.Request.Context(), userID, targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat_id": chatID,
		"message": "chat ready",
	})
}
