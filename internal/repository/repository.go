package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	ListActiveByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error)
	Update(ctx context.Context, user models.User) error
	ListReviewPRs(ctx context.Context, userID string) ([]models.PullRequest, error)
	WithTx(tx *gorm.DB) UserRepository
}

type TeamRepository interface {
	Create(ctx context.Context, team models.Team) error
	GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	GetByName(ctx context.Context, teamName string) (*models.Team, error)

	ListUsersByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error)
	WithTx(tx *gorm.DB) TeamRepository
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr models.PullRequest) error
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
	Update(ctx context.Context, pr models.PullRequest) error

	AddReviewer(ctx context.Context, prID string, reviewerID string) error
	RemoveReviewer(ctx context.Context, prID string, reviewerID string) error

	ListReviewers(ctx context.Context, prID string) ([]models.User, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequest, error)
	WithTx(tx *gorm.DB) PullRequestRepository
	RemoveReviewerFromAllPRs(ctx context.Context, userID string) error
}

type ReviewerHistoryRepository interface {
	AddEvent(ctx context.Context, event models.ReviewerAssignmentHistory) error
	CountAssignmentsByUsers(ctx context.Context) ([]dto.ReviewerStatsItem, error)
	WithTx(tx *gorm.DB) ReviewerHistoryRepository
}
