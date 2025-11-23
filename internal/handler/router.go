package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/mink0ff/pr_service/internal/service"
)

func RegisterRoutes(r chi.Router, ts service.TeamService, us service.UserService, prs service.PRService, ss service.StatsService) {
	teamHandler := NewTeamHandler(ts)
	r.Post("/team/add", teamHandler.CreateTeam)
	r.Get("/team/get", teamHandler.GetTeam)

	userHandler := NewUserHandler(us)
	r.Post("/users/setIsActive", userHandler.SetActive)
	r.Get("/users/getReview", userHandler.GetReviewPRs)

	prHandler := NewPRHandler(prs)
	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Post("/pullRequest/merge", prHandler.MergePR)
	r.Post("/pullRequest/reassign", prHandler.ReassignReviewer)

	statsHandler := NewStatsHandler(ss)
	r.Get("/stats/reviewers", statsHandler.GetReviewerStatsHandler)
}
