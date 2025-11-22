package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).First(&user, "user_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

func (r *UserRepo) ListActiveByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Where("team_id = ? AND is_active = TRUE", teamID).
		Find(&users).Error

	return users, err
}

func (r *UserRepo) ListReviewPRs(ctx context.Context, userID uuid.UUID) ([]models.PullRequest, error) {
	var PullRequest []models.PullRequest
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers pr ON pr.pull_request_id = pull_requests.pull_request_id").
		Where("pr.reviewer_id = ?", userID).
		Find(&PullRequest).Error

	return PullRequest, err
}
