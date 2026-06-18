package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(pr *repository.PostRepository) *PostService {
	return &PostService{postRepo: pr}
}

func (s *PostService) CreatePost(ctx context.Context, userID uuid.UUID, content, mediaURL string) (*models.Post, error) {
	post := &models.Post{
		ID:        uuid.New(),
		AuthorID:  userID,
		Content:   content,
		MediaURL:  mediaURL,
		CreatedAt: time.Now(),
	}
	if err := s.postRepo.CreatePost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetGlobalFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	return s.postRepo.GetFeed(ctx, limit, offset, userID)
}
