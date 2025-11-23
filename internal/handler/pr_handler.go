package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/service"
)

type PRHandler struct {
	prService service.PRService
}

func NewPRHandler(prService service.PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.prService.CreatePR(r.Context(), &req)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.prService.MergePR(r.Context(), &req)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.prService.ReassignReviewer(r.Context(), &req)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
