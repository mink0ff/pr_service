package handler

import (
	"net/http"

	"github.com/mink0ff/pr_service/internal/service"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetReviewerStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.statsService.GetReviewerStats(ctx)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
