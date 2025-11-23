package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.teamService.CreateTeam(r.Context(), &req)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := chi.URLParam(r, "team_name")
	if teamName == "" {
		teamName = r.URL.Query().Get("team_name")
		if teamName == "" {
			http.Error(w, "team_name query is required", http.StatusBadRequest)
			return
		}
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, team)
}
