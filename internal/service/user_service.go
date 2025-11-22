package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	userID := uuid.New()
	user := models.User{
		UserID:   userID,
		Username: req.Name,
		TeamID:   req.TeamID,
		IsActive: req.IsActive,
	}

	err := s.userRepo.Create(ctx, user)

	if err != nil {
		log.Printf("user create error: %v", err)
		return nil, ErrUserExists
	}

	return &user, nil
}

func (s *UserService) SetActive(ctx context.Context, req dto.SetUserActiveRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil || user == nil {
		log.Printf("user not found: %v", req.UserID)
		return nil, ErrUserNotFound
	}

	user.IsActive = req.IsActive

	if err := s.userRepo.Update(ctx, *user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetReviewPRs(ctx context.Context, userID uuid.UUID) ([]models.PullRequest, error) {
	return s.userRepo.ListReviewPRs(ctx, userID)
}
