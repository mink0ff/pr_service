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

func (r *TeamRepo) AddUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	tm := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
	}
	return r.db.WithContext(ctx).Create(&tm).Error
}

func (r *TeamRepo) RemoveUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&models.TeamMember{}).Error
}

func (r *TeamRepo) ListUser(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Joins("JOIN team_members tm ON tm.user_id = users.user_id").
		Where("tm.team_id = ?", teamID).
		Find(&users).Error

	return users, err
}
