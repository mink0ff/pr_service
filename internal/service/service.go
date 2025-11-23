package service

import (
	"context"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.User, error)
	SetActive(ctx context.Context, req dto.SetUserActiveRequest) (*dto.User, error)
	GetReviewPRs(ctx context.Context, userID string) ([]models.PullRequest, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, req *dto.CreateTeamRequest) (*dto.CreateTeamResponse, error)
	GetTeam(ctx context.Context, teamName string) (*dto.Team, error)
}

type PRService interface {
	CreatePR(ctx context.Context, req *dto.CreatePRRequest) (*dto.CreatePRResponse, error)
	ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error)
	MergePR(ctx context.Context, req *dto.MergePRRequest) (*dto.MergePRResponse, error)
}
