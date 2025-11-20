package service

import (
	"context"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	id, err := s.repo.Create(ctx, models.User{
		Name:     req.Name,
		IsActive: req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       id,
		Name:     req.Name,
		IsActive: req.IsActive,
	}, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		IsActive: user.IsActive,
	}, nil
}
