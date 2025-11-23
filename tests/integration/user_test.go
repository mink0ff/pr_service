package integration

import (
	"context"
	"testing"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/tests/utils"
	"github.com/stretchr/testify/require"
)

func initUserServiceTest(t *testing.T) {
	utils.TruncateTables(ts.DB)
}

func TestUserService_SetActiveAndGetReviewPRs(t *testing.T) {
	initUserServiceTest(t)
	ctx := context.Background()

	team := dto.Team{
		TeamName: "frontend",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}
	_, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)

	user := dto.SetUserActiveRequest{UserID: "u1", IsActive: true}
	userResp, err := ts.UserService.SetActive(ctx, user)
	require.NoError(t, err)
	require.Equal(t, "u1", userResp.UserID)
	require.True(t, userResp.IsActive)

	userDeactivate := dto.SetUserActiveRequest{UserID: "u1", IsActive: false}
	userResp2, err := ts.UserService.SetActive(ctx, userDeactivate)
	require.NoError(t, err)
	require.False(t, userResp2.IsActive)

	nonExistent := dto.SetUserActiveRequest{UserID: "u999", IsActive: true}
	_, err = ts.UserService.SetActive(ctx, nonExistent)
	require.Error(t, err)
}

func TestUserService_GetReviewPRs(t *testing.T) {
	initUserServiceTest(t)
	ctx := context.Background()

	team := dto.Team{
		TeamName: "backend",
		Members: []dto.TeamMember{
			{UserID: "u3", Username: "Charlie", IsActive: true},
			{UserID: "u4", Username: "Dana", IsActive: true},
		},
	}
	createdTeam, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)
	require.Len(t, createdTeam.Team.Members, 2)

	fetchedTeam, err := ts.TeamService.GetTeam(ctx, "backend")
	require.NoError(t, err)
	require.Len(t, fetchedTeam.Members, 2)

	prReq := dto.CreatePRRequest{
		PullRequestID:   "pr-1001",
		PullRequestName: "pr-name",
		AuthorID:        "u3",
	}
	createdPR, err := ts.PRService.CreatePR(ctx, &prReq)
	require.NoError(t, err)
	require.Equal(t, "pr-1001", createdPR.PR.PullRequestID)

	reviewPRsU3, err := ts.UserService.GetReviewPRs(ctx, "u3")
	require.NoError(t, err)
	require.Len(t, reviewPRsU3, 0)

	reviewPRsU4, err := ts.UserService.GetReviewPRs(ctx, "u4")
	require.NoError(t, err)
	require.Len(t, reviewPRsU4, 1)
	require.Equal(t, "pr-1001", reviewPRsU4[0].PullRequestID)

	reviewPRs, err := ts.UserService.GetReviewPRs(ctx, "u999")
	require.NoError(t, err)
	require.Len(t, reviewPRs, 0)
}
