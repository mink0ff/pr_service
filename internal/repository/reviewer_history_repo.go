package repository

import (
	"context"
	"log"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type ReviewerHistoryRepo struct {
	db *gorm.DB
}

func NewReviewerHistoryRepo(db *gorm.DB) ReviewerHistoryRepository {
	return &ReviewerHistoryRepo{db: db}
}

func (r *ReviewerHistoryRepo) AddEvent(ctx context.Context, event models.ReviewerAssignmentHistory) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *ReviewerHistoryRepo) CountAssignmentsByUsers(ctx context.Context) ([]dto.ReviewerStatsItem, error) {
	var statsItems []dto.ReviewerStatsItem

	err := r.db.WithContext(ctx).
		Model(&models.ReviewerAssignmentHistory{}).
		Select("user_id, COUNT(*) AS count").
		Group("user_id").
		Order("count DESC").
		Scan(&statsItems).Error

	if err != nil {
		log.Printf("Failed to count reviewer assignments: %v\n", err)
		return nil, err
	}

	log.Printf("Successfully counted assignments for %d users\n", len(statsItems))
	return statsItems, nil
}

func (r *ReviewerHistoryRepo) WithTx(tx *gorm.DB) ReviewerHistoryRepository {
	return &ReviewerHistoryRepo{db: tx}
}
