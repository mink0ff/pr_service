package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type TeamRepo struct {
	db *gorm.DB
}

func NewTeamRepo(db *gorm.DB) TeamRepository {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team models.Team) error {
	return r.db.WithContext(ctx).Create(&team).Error
}

func (r *TeamRepo) GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).First(&team, "team_id = ?", teamID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &team, err
}

func (r *TeamRepo) GetByName(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).First(&team, "team_name = ?", teamName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &team, err
}

func (r *TeamRepo) ListUsersByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.User{}, nil
	}
	return users, err
}

func (r *TeamRepo) WithTx(tx *gorm.DB) TeamRepository {
	return &TeamRepo{db: tx}
}
