package services

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrPRMerged             = errors.New("pull request is already merged")
	ErrNoAvailableReviewers = errors.New("no available reviewers")
)
