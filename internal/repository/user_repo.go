package repository

import (
	"context"
	"errors"

	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (int, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *UserRepo) ListActiveByTeam(ctx context.Context, teamID int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Joins("JOIN team_users tu ON tu.user_id = user.id").
		Where("tu.team_id = ? AND users.is_active = true", teamID).Find(&users).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return users, err
}
