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

func NewTeamRepo(db *gorm.DB) *TeamRepo {
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

func (r *TeamRepo) AddUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("team_id", teamID).Error
}

func (r *TeamRepo) RemoveUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("user_id = ? AND team_id = ?", userID, teamID).
		Update("team_id", nil).Error
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
