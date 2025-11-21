package models

import "github.com/google/uuid"

type PRReviewer struct {
	PullRequestID uuid.UUID `db:"pull_request_id"`
	ReviewerID    uuid.UUID `db:"reviewer_id"`
}
