package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type RelationService struct {
	relRepo  *repository.RelationRepository
	postRepo *repository.PostRepository
}

func NewRelationService(relRepo *repository.RelationRepository, postRepo *repository.PostRepository) *RelationService {
	return &RelationService{relRepo: relRepo, postRepo: postRepo}
}

func (s *RelationService) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return errors.New("you cannot follow yourself") // Импортируй errors
	}
	return s.relRepo.Follow(ctx, followerID, followingID)
}

func (s *RelationService) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.relRepo.Unfollow(ctx, followerID, followingID)
}

func (s *RelationService) GetPersonalFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Post, error) {
	following, err := s.relRepo.GetFollowingIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.postRepo.GetFeedByAuthors(ctx, following, limit, offset)
}
