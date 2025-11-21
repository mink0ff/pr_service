package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListActiveByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error)
	Update(ctx context.Context, user models.User) error
	ListReviewPRs(ctx context.Context, userID uuid.UUID) ([]models.PullRequest, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team models.Team) error
	GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	GetByName(ctx context.Context, teamName string) (*models.Team, error)

	AddUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	RemoveUser(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error

	ListUser(ctx context.Context, teamID uuid.UUID) ([]models.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr models.PullRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PullRequest, error)
	Update(ctx context.Context, pr models.PullRequest) error

	AddReviewer(ctx context.Context, prID uuid.UUID, reviewerID uuid.UUID) error
	RemoveReviewer(ctx context.Context, prID uuid.UUID, reviewerID uuid.UUID) error

	ListReviewers(ctx context.Context, prID uuid.UUID) ([]models.User, error)
	ListByReviewer(ctx context.Context, reviewerID uuid.UUID) ([]models.PullRequest, error)
}
