package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(ur *repository.UserRepository) *UserService {
	return &UserService{userRepo: ur}
}

func (s *UserService) SearchUsers(ctx context.Context, query string) ([]*models.User, error) {
	return s.userRepo.SearchUsers(ctx, query)
}
