package integration

import (
	"context"
	"testing"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/tests/utils"
	"github.com/stretchr/testify/require"
)

func initStatsTest(t *testing.T) {
	utils.TruncateTables(ts.DB)
}

func TestStatsService_GetReviewerStats(t *testing.T) {
	initStatsTest(t)
	ctx := context.Background()

	stats, err := ts.StatsService.GetReviewerStats(ctx)
	require.NoError(t, err)
	require.Len(t, stats.Items, 0)

	team := dto.Team{
		TeamName: "analytics",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
			{UserID: "u3", Username: "Charlie", IsActive: true},
			{UserID: "u4", Username: "David", IsActive: true}, // для Reassign
		},
	}
	_, err = ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)

	prReq := dto.CreatePRRequest{
		PullRequestID:   "pr-1001",
		PullRequestName: "Add analytics module",
		AuthorID:        "u1",
	}
	createdPR, err := ts.PRService.CreatePR(ctx, &prReq)
	require.NoError(t, err)

	stats, err = ts.StatsService.GetReviewerStats(ctx)
	require.NoError(t, err)
	require.Len(t, stats.Items, 2)

	assigned := createdPR.PR.AssignedReviewers
	for _, item := range stats.Items {
		require.Contains(t, assigned, item.UserID)
		require.Equal(t, int64(1), item.Count)
	}
}
