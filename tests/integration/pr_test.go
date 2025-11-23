package integration

import (
	"context"
	"testing"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/service"
	"github.com/mink0ff/pr_service/tests/utils"
	"github.com/stretchr/testify/require"
)

func initPRServiceTest(t *testing.T) {
	utils.TruncateTables(ts.DB)
}

func TestPRService_FullFlow(t *testing.T) {
	initPRServiceTest(t)
	ctx := context.Background()

	team := dto.Team{
		TeamName: "backend",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
			{UserID: "u3", Username: "Charlie", IsActive: true},
			{UserID: "u4", Username: "Alex", IsActive: true},
		},
	}
	_, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)

	prReq := dto.CreatePRRequest{
		PullRequestID:   "pr-1001",
		PullRequestName: "Add search",
		AuthorID:        "u1",
	}
	createdPR, err := ts.PRService.CreatePR(ctx, &prReq)
	require.NoError(t, err)
	require.Equal(t, "pr-1001", createdPR.PR.PullRequestID)
	require.Equal(t, dto.PRStatusOpen, createdPR.PR.Status)
	require.Len(t, createdPR.PR.AssignedReviewers, 2)

	oldReviewer := createdPR.PR.AssignedReviewers[0]
	reassignReq := dto.ReassignReviewerRequest{
		PullRequestID: "pr-1001",
		OldUserID:     oldReviewer,
	}
	reassignedPR, err := ts.PRService.ReassignReviewer(ctx, &reassignReq)
	require.NoError(t, err)
	require.Equal(t, "pr-1001", reassignedPR.PR.PullRequestID)
	require.NotEqual(t, oldReviewer, reassignedPR.ReplacedBy)
	require.Contains(t, reassignedPR.PR.AssignedReviewers, reassignedPR.ReplacedBy)

	mergeReq := dto.MergePRRequest{PullRequestID: "pr-1001"}
	mergedPR, err := ts.PRService.MergePR(ctx, &mergeReq)
	require.NoError(t, err)
	require.Equal(t, dto.PRStatusMerged, mergedPR.PR.Status)

	mergedPR2, err := ts.PRService.MergePR(ctx, &mergeReq)
	require.NoError(t, err)
	require.Equal(t, dto.PRStatusMerged, mergedPR2.PR.Status)
}

func TestPRService_CreatePRWithNonExistentAuthor(t *testing.T) {
	initPRServiceTest(t)
	ctx := context.Background()

	prReq := dto.CreatePRRequest{
		PullRequestID:   "pr-2001",
		PullRequestName: "Add feature X",
		AuthorID:        "u999",
	}
	_, err := ts.PRService.CreatePR(ctx, &prReq)
	require.Error(t, err)
}

func TestPRService_MergeAndReassignNonExistentPR(t *testing.T) {
	initPRServiceTest(t)
	ctx := context.Background()

	mergeReq := dto.MergePRRequest{PullRequestID: "pr-999"}
	_, err := ts.PRService.MergePR(ctx, &mergeReq)
	require.Error(t, err)

	reassignReq := dto.ReassignReviewerRequest{
		PullRequestID: "pr-999",
		OldUserID:     "u1",
	}
	_, err = ts.PRService.ReassignReviewer(ctx, &reassignReq)
	require.Error(t, err)
}

func TestPRService_ReassignAfterMerge(t *testing.T) {
	initPRServiceTest(t)
	ctx := context.Background()

	team := dto.Team{
		TeamName: "frontend",
		Members: []dto.TeamMember{
			{UserID: "u10", Username: "Alice", IsActive: true},
			{UserID: "u11", Username: "Bob", IsActive: true},
		},
	}
	_, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)

	prReq := dto.CreatePRRequest{
		PullRequestID:   "pr-3001",
		PullRequestName: "Feature Y",
		AuthorID:        "u10",
	}
	createdPR, err := ts.PRService.CreatePR(ctx, &prReq)
	require.NoError(t, err)

	mergeReq := dto.MergePRRequest{PullRequestID: "pr-3001"}
	_, err = ts.PRService.MergePR(ctx, &mergeReq)
	require.NoError(t, err)

	reassignReq := dto.ReassignReviewerRequest{
		PullRequestID: "pr-3001",
		OldUserID:     createdPR.PR.AssignedReviewers[0],
	}
	_, err = ts.PRService.ReassignReviewer(ctx, &reassignReq)
	require.ErrorIs(t, err, service.ErrPRMerged)
}
