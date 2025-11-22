package handler

import (
	"errors"
	"net/http"

	svc "github.com/mink0ff/pr_service/internal/service"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func MapError(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, svc.ErrTeamExists):
		return http.StatusConflict, ErrorResponse{Code: "TEAM_EXISTS", Message: err.Error()}

	case errors.Is(err, svc.ErrTeamNotFound):
		return http.StatusNotFound, ErrorResponse{Code: "TEAM_NOT_FOUND", Message: err.Error()}

	case errors.Is(err, svc.ErrUserNotFound):
		return http.StatusNotFound, ErrorResponse{Code: "USER_NOT_FOUND", Message: err.Error()}

	case errors.Is(err, svc.ErrPRExists):
		return http.StatusConflict, ErrorResponse{Code: "PR_EXISTS", Message: err.Error()}

	case errors.Is(err, svc.ErrPRNotFound):
		return http.StatusNotFound, ErrorResponse{Code: "PR_NOT_FOUND", Message: err.Error()}

	case errors.Is(err, svc.ErrPRMerged):
		return http.StatusConflict, ErrorResponse{Code: "PR_MERGED", Message: err.Error()}

	case errors.Is(err, svc.ErrReviewerNotAssigned):
		return http.StatusBadRequest, ErrorResponse{Code: "NOT_ASSIGNED", Message: err.Error()}

	case errors.Is(err, svc.ErrNoCandidate):
		return http.StatusConflict, ErrorResponse{Code: "NO_CANDIDATE", Message: err.Error()}
	}

	return http.StatusInternalServerError,
		ErrorResponse{Code: "INTERNAL_ERROR", Message: "internal server error"}
}
