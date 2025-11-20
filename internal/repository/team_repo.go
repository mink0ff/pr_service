package repository

import (
	"context"
	"errors"

	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type TeamRepo struct {
	db *gorm.DB
}

func NewTeamRepo(db *gorm.DB) *TeamRepo {
	return &TeamRepo{}
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) (int, error) {
	if err := r.db.Create(team).Error; err != nil {
		return 0, err
	}

	return team.ID, nil
}

func (r *TeamRepo) GetByID(ctx context.Context, id int) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).First(&team, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &team, nil
}

func (r *TeamRepo) AddUser(ctx context.Context, teamID, userID int) error {
	tu := models.TeamUser{TeamID: teamID, UserID: userID}
	return r.db.WithContext(ctx).Create(&tu).Error
}

func (r *TeamRepo) RemoveUser(ctx context.Context, teamID, userID int) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&models.TeamUser{}).Error
}

func (r *TeamRepo) ListUsers(ctx context.Context, teamID int) ([]models.TeamUser, error) {
	var users []models.TeamUser
	err := r.db.WithContext(ctx).
		Joins("JOIN team_users tu ON tu.user_id = users.id").
		Where("tu.team_id = ?", teamID).
		Find(&users).Error

	return users, err
}
