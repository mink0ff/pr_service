package models

import "github.com/google/uuid"

type User struct {
	UserID   string    `db:"user_id"`
	Username string    `db:"username"`
	TeamID   uuid.UUID `db:"team_id"`
	IsActive bool      `db:"is_active"`
}
