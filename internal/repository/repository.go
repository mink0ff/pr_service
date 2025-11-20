package repository

import (
	"context"
	_ "context"

	"github.com/mink0ff/pr_service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	ListActiveByTeam(ctx context.Context, teamID int) ([]models.User, error)
	Update(ctx context.Context, user models.User) error
}

type TeamRepository interface {
	Create(ctx context.Context) (int, error)
	GetByID(ctx context.Context) (models.Team, error)
	AddUser(ctx context.Context) error
	RemoveUser(ctx context.Context) error
	ListUser(ctx context.Context) ([]models.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context) (int, error)
	GetByID(ctx context.Context) (*models.PullRequest, error)
	Update(ctx context.Context) error
	AddReviewer(ctx context.Context) error
	RemoveReviewer(ctx context.Context) error
	ListReviewers(ctx context.Context) ([]models.User, error)
	ListByReviewer(ctx context.Context) ([]models.PullRequest, error)
}
