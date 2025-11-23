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
	txManager *transaction.Manager
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, manager *transaction.Manager) TeamService {
	return &TeamServiceImpl{teamRepo: teamRepo, userRepo: userRepo, txManager: manager}
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
