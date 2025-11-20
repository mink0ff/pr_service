package repository

import (
	"context"
	"errors"

	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type PrRepo struct {
	db *gorm.DB
}

func NewPrRepo(db *gorm.DB) *PrRepo {
	return &PrRepo{db: db}
}

func (r *PrRepo) Create(ctx context.Context, pr *models.PullRequest) (int, error) {
	if err := r.db.WithContext(ctx).Create(&pr).Error; err != nil {
		return 0, err
	}

	return pr.ID, nil
}

func (r *PrRepo) GetByID(ctx context.Context, id int) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.WithContext(ctx).First(&pr, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &pr, err
}

func (r *PrRepo) Update(ctx context.Context, pr *models.PullRequest) error {
	return r.db.WithContext(ctx).Save(&pr).Error
}

func (r *PrRepo) AddReviewer(ctx context.Context, prID int, userID int) error {
	return r.db.WithContext(ctx).Create(&models.PRReviewer{ID: prID, UserID: userID}).Error
}

func (r *PrRepo) DeleteReviewer(ctx context.Context, prID int, userID int) error {
	return r.db.WithContext(ctx).
		Where("pr_id = ? AND user_id = ?", prID, userID).
		Delete(&models.PRReviewer{}).Error
}

func (r *PrRepo) ListReviewers(ctx context.Context, prID int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewer prr ON prr.user_id = users.id").
		Where("prr.pr_id = ?", prID).
		Find(&users).Error

	return users, err
}

func (r *PrRepo) ListByReviewer(ctx context.Context, userID int) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewer prr ON prr.pr_id = pull_request.id").
		Where("prr.user_id = ?", userID).
		Find(&prs).Error

	return prs, err
}
