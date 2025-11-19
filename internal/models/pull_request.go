package models

import (
	"time"
)

type RPStatus string

const (
	RPOpen   RPStatus = "OPEN"
	RPMerged RPStatus = "MERGED"
)

type PullRequest struct {
	ID        int
	Title     string
	AuthorID  int
	TeamID    int
	Status    RPStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PRReviewer struct {
	ID         int
	UserID     int
	AssignedAt time.Time
}
