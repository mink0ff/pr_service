package models

import (
	"time"

	"github.com/google/uuid"
)

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   uuid.UUID  `db:"pull_request_id"`
	PullRequestName string     `db:"pull_request_name"`
	AuthorID        uuid.UUID  `db:"author_id"`
	Status          PRStatus   `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	MergedAt        *time.Time `db:"merged_at"`
}
