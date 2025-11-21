package models

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	TeamID   uuid.UUID `db:"team_id"`
	IsActive bool      `db:"is_active"`
}
