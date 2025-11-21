package models

import "github.com/google/uuid"

type Team struct {
	TeamID   uuid.UUID `db:"team_id"`
	TeamName string    `db:"team_name"`
}
