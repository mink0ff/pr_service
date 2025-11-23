package service

import (
	"context"
	_ "log"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
	"github.com/mink0ff/pr_service/internal/repository/transaction"
	"gorm.io/gorm"
)

type TeamServiceImpl struct {
	teamRepo  repository.TeamRepository
	userRepo  repository.UserRepository
	prRepo    repository.PullRequestRepository
	txManager *transaction.Manager
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, prRepo repository.PullRequestRepository, manager *transaction.Manager) TeamService {
	return &TeamServiceImpl{teamRepo: teamRepo, userRepo: userRepo, prRepo: prRepo, txManager: manager}
}

func (s *TeamServiceImpl) CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateTeamResponse, error) {
	var resp *dto.CreateTeamResponse

	err := s.txManager.Do(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		txTeamRepo := s.teamRepo.WithTx(tx)
		txUserRepo := s.userRepo.WithTx(tx)

		team, err := s.createTeam(txCtx, req, txTeamRepo)
		if err != nil {
			return err
		}

		if err := s.createOrUpdateMembers(txCtx, team.TeamID, req.Members, txUserRepo); err != nil {
			return err
		}

		resp = &dto.CreateTeamResponse{
			Team: dto.Team{
				TeamName: team.TeamName,
				Members:  req.Members,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *TeamServiceImpl) GetTeam(ctx context.Context, teamName string) (*dto.Team, error) {
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
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		}
	}

	return &dto.Team{
		TeamName: team.TeamName,
		Members:  members,
	}, nil
}

func (s *TeamServiceImpl) createTeam(
	ctx context.Context,
	req *dto.CreateTeamRequest,
	teamRepo repository.TeamRepository,
) (*models.Team, error) {

	existing, err := teamRepo.GetByName(ctx, req.TeamName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTeamExists
	}

	team := &models.Team{
		TeamID:   uuid.New(),
		TeamName: req.TeamName,
	}

	if err := teamRepo.Create(ctx, *team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamServiceImpl) createOrUpdateMembers(
	ctx context.Context,
	teamID uuid.UUID,
	members []dto.TeamMember,
	userRepo repository.UserRepository,
) error {

	for _, m := range members {
		existingUser, err := userRepo.GetByID(ctx, m.UserID)
		if err != nil {
			return err
		}

		user := models.User{
			UserID:   m.UserID,
			Username: m.Username,
			TeamID:   teamID,
			IsActive: m.IsActive,
		}

		if existingUser == nil {
			if err := userRepo.Create(ctx, user); err != nil {
				return err
			}
		} else {
			if err := userRepo.Update(ctx, user); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TeamServiceImpl) DeactivateTeamUsers(ctx context.Context, req *dto.DeactivateTeamUsersRequest) (*dto.DeactivateTeamUsersResponse, error) {
	var resp *dto.DeactivateTeamUsersResponse

	err := s.txManager.Do(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		txUserRepo := s.userRepo.WithTx(tx)
		txPrRepo := s.prRepo.WithTx(tx)

		team, err := s.teamRepo.GetByName(txCtx, req.TeamName)
		if err != nil {
			return err
		}
		if team == nil {
			return ErrTeamNotFound
		}

		users, err := txUserRepo.ListActiveByTeam(txCtx, team.TeamID)
		if err != nil {
			return err
		}

		for _, u := range users {
			u.IsActive = false
			if err := txUserRepo.Update(txCtx, u); err != nil {
				return err
			}

			if err := txPrRepo.RemoveReviewerFromAllPRs(txCtx, u.UserID); err != nil {
				return err
			}
		}

		resp = &dto.DeactivateTeamUsersResponse{
			TeamName:         req.TeamName,
			DeactivatedCount: len(users),
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
