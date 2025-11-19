package repository

import (
	"context"
	"github.com/mink0ff/pr_service/internal/models"
)

type UserRepository interface {
	Create() (int, error)
	GetByID() (*models.User, error)
	ListActiveByTeam() ([]models.User, error)
	Update() error
}

type TeamRepository interface {
	Create() (int, error)
	GetByID() (models.Team, error)
	AddUser() error
	RemoveUser() error
	ListUser() ([]models.User, error)
}

type PullRequestRepository interface {
	Create() (int, error)
	GetByID() (*models.PullRequest, error)
	Update() error
	AddReviewer() error
	RemoveReviewer() error
	ListReviewers() ([]models.User, error)
	ListByReviewer() ([]models.PullRequest, error)
}
