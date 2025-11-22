package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// POST /users/setIsActive
func (h *UserHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.SetActive(r.Context(), req)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

// GET /users/getReview?user_id=xxx
func (h *UserHandler) GetReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query is required", http.StatusBadRequest)
		return
	}

	prs, err := h.userService.GetReviewPRs(r.Context(), userID)
	if err != nil {
		status, errResp := MapError(err)
		writeJSON(w, status, errResp)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
