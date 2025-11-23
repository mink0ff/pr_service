package repository

import (
	"context"
	"errors"

	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PrRepo struct {
	db *gorm.DB
}

func NewPrRepo(db *gorm.DB) PullRequestRepository {
	return &PrRepo{db: db}
}

func (r *PrRepo) Create(ctx context.Context, pr models.PullRequest) error {
	return r.db.WithContext(ctx).Create(&pr).Error
}

func (r *PrRepo) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&pr, "pull_request_id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pr, err
}

func (r *PrRepo) Update(ctx context.Context, pr models.PullRequest) error {
	return r.db.WithContext(ctx).
		Where("pull_request_id = ?", pr.PullRequestID).
		Save(&pr).Error
}

func (r *PrRepo) AddReviewer(ctx context.Context, prID string, reviewerID string) error {
	record := models.PRReviewer{
		PullRequestID: prID,
		ReviewerID:    reviewerID,
	}
	return r.db.WithContext(ctx).Create(&record).Error
}

func (r *PrRepo) RemoveReviewer(ctx context.Context, prID string, reviewerID string) error {
	return r.db.WithContext(ctx).
		Where("pull_request_id = ? AND reviewer_id = ?", prID, reviewerID).
		Delete(&models.PRReviewer{}).Error
}

func (r *PrRepo) ListReviewers(ctx context.Context, prID string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers prr ON prr.reviewer_id = users.user_id").
		Where("prr.pull_request_id = ?", prID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&users).Error
	return users, err
}

func (r *PrRepo) ListByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest

	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.reviewer_id = ?", reviewerID).
		Find(&prs).Error

	return prs, err
}

func (r *PrRepo) WithTx(tx *gorm.DB) PullRequestRepository {
	return &PrRepo{db: tx}
}

func (r *PrRepo) RemoveReviewerFromAllPRs(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("reviewer_id = ?", userID).
		Delete(&models.PRReviewer{}).Error
}
