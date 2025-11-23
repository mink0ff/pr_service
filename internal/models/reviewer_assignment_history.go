package models

import (
	"time"

	"github.com/google/uuid"
)

type ReviewerAssignmentHistory struct {
	AssigmentHistoryID uuid.UUID `db:"assigment_history_id"`
	PrID               string    `db:"pr_id"`
	UserID             string    `db:"user_id"`
	CreatedAt          time.Time `db:"created_at"`
}
