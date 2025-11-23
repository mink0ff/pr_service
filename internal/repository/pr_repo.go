package repository

import (
	"context"
	"errors"
	"log"

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
	err := r.db.WithContext(ctx).Create(&pr).Error
	if err != nil {
		log.Printf("Failed to create PullRequest %v: %v\n", pr.PullRequestID, err)
	} else {
		log.Printf("PullRequest %v created successfully\n", pr.PullRequestID)
	}
	return err
}

func (r *PrRepo) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	var pr models.PullRequest
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&pr, "pull_request_id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("PullRequest %v not found\n", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error fetching PullRequest %v: %v\n", id, err)
	} else {
		log.Printf("PullRequest %v fetched successfully\n", id)
	}
	return &pr, err
}

func (r *PrRepo) Update(ctx context.Context, pr models.PullRequest) error {
	err := r.db.WithContext(ctx).
		Where("pull_request_id = ?", pr.PullRequestID).
		Save(&pr).Error
	if err != nil {
		log.Printf("Failed to update PullRequest %v: %v\n", pr.PullRequestID, err)
	} else {
		log.Printf("PullRequest %v updated successfully\n", pr.PullRequestID)
	}
	return err
}

func (r *PrRepo) AddReviewer(ctx context.Context, prID string, reviewerID string) error {
	record := models.PRReviewer{
		PullRequestID: prID,
		ReviewerID:    reviewerID,
	}
	err := r.db.WithContext(ctx).Create(&record).Error
	if err != nil {
		log.Printf("Failed to add reviewer %v to PR %v: %v\n", reviewerID, prID, err)
	} else {
		log.Printf("Reviewer %v added to PR %v successfully\n", reviewerID, prID)
	}
	return err
}

func (r *PrRepo) RemoveReviewer(ctx context.Context, prID string, reviewerID string) error {
	err := r.db.WithContext(ctx).
		Where("pull_request_id = ? AND reviewer_id = ?", prID, reviewerID).
		Delete(&models.PRReviewer{}).Error
	if err != nil {
		log.Printf("Failed to remove reviewer %v from PR %v: %v\n", reviewerID, prID, err)
	} else {
		log.Printf("Reviewer %v removed from PR %v successfully\n", reviewerID, prID)
	}
	return err
}

func (r *PrRepo) ListReviewers(ctx context.Context, prID string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers prr ON prr.reviewer_id = users.user_id").
		Where("prr.pull_request_id = ?", prID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&users).Error
	if err != nil {
		log.Printf("Failed to list reviewers for PR %v: %v\n", prID, err)
	} else {
		log.Printf("Found %d reviewers for PR %v\n", len(users), prID)
	}
	return users, err
}

func (r *PrRepo) ListByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest

	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers prr ON prr.pull_request_id = pull_requests.pull_request_id").
		Where("prr.reviewer_id = ?", reviewerID).
		Find(&prs).Error

	if err != nil {
		log.Printf("Failed to list PRs for reviewer %v: %v\n", reviewerID, err)
	} else {
		log.Printf("Found %d PRs for reviewer %v\n", len(prs), reviewerID)
	}
	return prs, err
}

func (r *PrRepo) WithTx(tx *gorm.DB) PullRequestRepository {
	return &PrRepo{db: tx}
}

func (r *PrRepo) RemoveReviewerFromAllPRs(ctx context.Context, userID string) error {
	err := r.db.WithContext(ctx).
		Where("reviewer_id = ?", userID).
		Delete(&models.PRReviewer{}).Error
	if err != nil {
		log.Printf("Failed to remove reviewer %v from all PRs: %v\n", userID, err)
	} else {
		log.Printf("Reviewer %v removed from all PRs successfully\n", userID)
	}
	return err
}
