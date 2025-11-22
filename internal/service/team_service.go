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

func (s *TeamService) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateTeamResponse, error) {
	existingTeam, err := s.teamRepo.GetByName(ctx, req.TeamName)
	if err != nil {
		log.Printf("team repo error: %v", err)
		return nil, err
	}
	if existingTeam != nil {
		return nil, ErrTeamExists
	}

	teamID := uuid.New()
	team := models.Team{
		TeamID:   teamID,
		TeamName: req.TeamName,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		log.Printf("team create error: %v", err)
		return nil, err
	}

	for _, member := range req.Members {
		userID, err := uuid.Parse(member.UserID)
		if err != nil {
			log.Printf("invalid user_id %s: %v", member.UserID, err)
			return nil, err
		}

		existingUser, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			log.Printf("user repo error: %v", err)
			return nil, err
		}

		if existingUser == nil {
			user := models.User{
				UserID:   userID,
				Username: member.Username,
				TeamID:   teamID,
				IsActive: member.IsActive,
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				log.Printf("create user error: %v", err)
				return nil, err
			}
		} else {
			existingUser.Username = member.Username
			existingUser.IsActive = member.IsActive
			existingUser.TeamID = teamID
			if err := s.userRepo.Update(ctx, *existingUser); err != nil {
				log.Printf("update user error: %v", err)
				return nil, err
			}
		}

		if err := s.teamRepo.AddUser(ctx, teamID, userID); err != nil {
			log.Printf("add user to team error: %v", err)
			return nil, err
		}
	}

	resp := &dto.CreateTeamResponse{
		Team: *req,
	}
	return resp, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*dto.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil || team == nil {
		return nil, ErrTeamNotFound
	}

	users, err := s.teamRepo.ListUsersByTeam(ctx, team.TeamID)
	if err != nil {
		return nil, err
	}

	members := make([]dto.TeamMember, len(users))
	for i, u := range users {
		members[i] = dto.TeamMember{
			UserID:   u.UserID.String(),
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	return &dto.Team{
		TeamName: team.TeamName,
		Members:  members,
	}, nil
}
