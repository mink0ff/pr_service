package integration

import (
	"context"
	"testing"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/tests/utils"
	"github.com/stretchr/testify/require"
)

func initTeamServiceTest(t *testing.T) {
	utils.TruncateTables(ts.DB)
}

func TestTeamService_CreateGetAndDeactivateTeam(t *testing.T) {
	initTeamServiceTest(t)

	ctx := context.Background()

	team := dto.Team{
		TeamName: "frontend",
		Members: []dto.TeamMember{
			{UserID: "u3", Username: "Alice", IsActive: true},
			{UserID: "u4", Username: "Bob1", IsActive: true},
			{UserID: "u5", Username: "Bob2", IsActive: true},
			{UserID: "u6", Username: "Bob3", IsActive: true},
		},
	}

	createdTeam, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)
	require.Equal(t, "frontend", createdTeam.Team.TeamName)
	require.Len(t, createdTeam.Team.Members, 4)

	fetchedTeam, err := ts.TeamService.GetTeam(ctx, "frontend")
	require.NoError(t, err)
	require.Equal(t, "frontend", fetchedTeam.TeamName)
	require.Len(t, fetchedTeam.Members, 4)

	for _, member := range fetchedTeam.Members {
		require.True(t, member.IsActive)
	}

	deactivateReq := dto.DeactivateTeamUsersRequest{
		TeamName: "frontend",
	}

	deactivateResp, err := ts.TeamService.DeactivateTeamUsers(ctx, &deactivateReq)
	require.NoError(t, err)
	require.Equal(t, "frontend", deactivateResp.TeamName)
	require.Equal(t, 4, deactivateResp.DeactivatedCount)

	fetchedTeamAfter, err := ts.TeamService.GetTeam(ctx, "frontend")
	require.NoError(t, err)
	for _, member := range fetchedTeamAfter.Members {
		require.False(t, member.IsActive)
	}
}

func TestTeamService_GetNonExistentTeam(t *testing.T) {
	initTeamServiceTest(t)
	ctx := context.Background()

	_, err := ts.TeamService.GetTeam(ctx, "nonexistent")
	require.Error(t, err)
}

func TestTeamService_DeactivateNonExistentTeam(t *testing.T) {
	initTeamServiceTest(t)
	ctx := context.Background()

	req := dto.DeactivateTeamUsersRequest{TeamName: "nonexistent"}
	_, err := ts.TeamService.DeactivateTeamUsers(ctx, &req)
	require.Error(t, err)
}

func TestTeamService_CreateTeamWithExistingUsers(t *testing.T) {
	initTeamServiceTest(t)
	ctx := context.Background()

	team := dto.Team{
		TeamName: "backend",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	createdTeam, err := ts.TeamService.CreateTeam(ctx, &team)
	require.NoError(t, err)
	require.Len(t, createdTeam.Team.Members, 2)

	newTeam := dto.Team{
		TeamName: "backend-renamed",
		Members: []dto.TeamMember{
			{UserID: "u1", Username: "Alice Updated", IsActive: true},
			{UserID: "u2", Username: "Bob Updated", IsActive: true},
		},
	}
	createdTeam2, err := ts.TeamService.CreateTeam(ctx, &newTeam)
	require.NoError(t, err)
	require.Len(t, createdTeam2.Team.Members, 2)

	for _, m := range createdTeam2.Team.Members {
		if m.UserID == "u1" {
			require.Equal(t, "Alice Updated", m.Username)
		}
		if m.UserID == "u2" {
			require.Equal(t, "Bob Updated", m.Username)
		}
	}
}
