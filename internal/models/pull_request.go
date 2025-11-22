package models

import (
	"time"
)

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string     `db:"pull_request_id"`
	PullRequestName string     `db:"pull_request_name"`
	AuthorID        string     `db:"author_id"`
	Status          PRStatus   `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	MergedAt        *time.Time `db:"merged_at"`
}
