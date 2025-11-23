package service

import (
	"context"
	"log"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(userRepo repository.UserRepository, teamRepo repository.TeamRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo, teamRepo: teamRepo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.User, error) {
	user := models.User{
		UserID:   req.UserID,
		Username: req.Name,
		TeamID:   req.TeamID,
		IsActive: req.IsActive,
	}

	err := s.userRepo.Create(ctx, user)

	if err != nil {
		log.Printf("user create error: %v", err)
		return nil, ErrUserExists
	}

	team, err := s.teamRepo.GetByID(ctx, user.TeamID)

	if err != nil {
		log.Printf("team get error: %v", err)
		return nil, ErrTeamNotFound
	}

	dtoUser := dto.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamName: team.TeamName,
		IsActive: user.IsActive,
	}

	return &dtoUser, nil
}

func (s *UserServiceImpl) SetActive(ctx context.Context, req dto.SetUserActiveRequest) (*dto.User, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil || user == nil {
		log.Printf("user not found: %v", req.UserID)
		return nil, ErrUserNotFound
	}

	userUpdate := models.User{
		UserID:   user.UserID,
		Username: user.Username,
		TeamID:   user.TeamID,
		IsActive: req.IsActive,
	}

	if err := s.userRepo.Update(ctx, userUpdate); err != nil {
		return nil, err
	}

	team, err := s.teamRepo.GetByID(ctx, user.TeamID)

	if err != nil {
		log.Printf("team get error: %v", err)
		return nil, ErrTeamNotFound
	}

	userDto := dto.User{
		UserID:   userUpdate.UserID,
		Username: userUpdate.Username,
		TeamName: team.TeamName,
		IsActive: userUpdate.IsActive,
	}

	return &userDto, nil
}

func (s *UserServiceImpl) GetReviewPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	return s.userRepo.ListReviewPRs(ctx, userID)
}
