package service

import (
	"context"

	"github.com/mink0ff/pr_service/internal/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUser(ctx context.Context, id int) (*dto.UserResponse, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, req dto.CreateTeamRequest) (*dto.TeamResponse, error)
	AddUser(ctx context.Context, teamID int, userID int) error
	GetTeam(ctx context.Context, id int) (*dto.TeamResponse, error)
}

type PRService interface {
	CreatePR(ctx context.Context, req dto.CreatePRRequest) (*dto.PRResponse, error)
	MergePR(ctx context.Context, prID int) (*dto.PRResponse, error)
	ReassignReviewer(ctx context.Context, prID int, reviewerID int) (*dto.PRResponse, error)
	GetPR(ctx context.Context, id int) (*dto.PRResponse, error)
	ListByReviewer(ctx context.Context, userID int) ([]dto.PRResponse, error)
}
