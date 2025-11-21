package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

func (s *TeamService) CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*models.Team, error) {
	teamID := uuid.New()

	team := models.Team{
		TeamID:   teamID,
		TeamName: req.TeamName,
	}

	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		log.Printf("team create error: %v", err)
		return nil, ErrTeamExists
	}

	for _, m := range req.Members {
		user := models.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamID:   team.TeamID,
			IsActive: m.IsActive,
		}

		err := s.userRepo.Create(ctx, user)
		if err != nil {
			if err := s.userRepo.Update(ctx, user); err != nil {
				log.Printf("team member update error: %v", err)
				return nil, err
			}
		}

		_ = s.teamRepo.AddUser(ctx, team.TeamID, user.UserID)
	}

	return &team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, []models.User, error) {
	team, err := s.teamRepo.GetByName(ctx, name)
	if err != nil || team == nil {
		return nil, nil, ErrTeamNotFound
	}

	users, err := s.teamRepo.ListUser(ctx, team.TeamID)
	if err != nil {
		return nil, nil, err
	}

	return team, users, nil
}
