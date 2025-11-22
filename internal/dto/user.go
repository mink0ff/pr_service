package dto

import "github.com/google/uuid"

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveResponse struct {
	User User `json:"user"`
}

type CreateUserRequest struct {
	UserID   string    `json:"user_id"`
	Name     string    `json:"username"`
	TeamID   uuid.UUID `json:"team_id"`
	IsActive bool      `json:"is_active"`
}

type GetReviewPRsResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}
