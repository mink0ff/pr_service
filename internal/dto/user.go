package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required"`
	TeamID   uuid.UUID `json:"team_id" binding:"required"`
	IsActive bool      `json:"isActive"`
}

type SetUserActiveRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	IsActive bool      `json:"is_active"`
}
type UserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}
