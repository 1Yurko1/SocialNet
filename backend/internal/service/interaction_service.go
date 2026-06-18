package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type InteractionService struct {
	interRepo *repository.InteractionRepository
}

func NewInteractionService(ir *repository.InteractionRepository) *InteractionService {
	return &InteractionService{interRepo: ir}
}

func (s *InteractionService) LikePost(ctx context.Context, userID, postID uuid.UUID) error {
	// Здесь можно добавить проверку: существует ли пост вообще
	return s.interRepo.LikePost(ctx, userID, postID)
}

func (s *InteractionService) UnlikePost(ctx context.Context, userID, postID uuid.UUID) error {
	return s.interRepo.UnlikePost(ctx, userID, postID)
}

func (s *InteractionService) AddComment(ctx context.Context, userID, postID uuid.UUID, content string) error {
	if content == "" {
		return errors.New("comment content cannot be empty")
	}
	return s.interRepo.AddComment(ctx, userID, postID, content)
}

func (s *InteractionService) GetComments(ctx context.Context, postID uuid.UUID) ([]*models.Comment, error) {
	return s.interRepo.GetCommentsByPost(ctx, postID)
}
