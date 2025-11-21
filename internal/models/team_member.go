package models

import "github.com/google/uuid"

type TeamMember struct {
	TeamID uuid.UUID `db:"team_id"`
	UserID uuid.UUID `db:"user_id"`
}
