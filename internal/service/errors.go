package service

import "errors"

var (
	ErrTeamExists          = errors.New("team already exists")
	ErrTeamNotFound        = errors.New("team not found")
	ErrUserExists          = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrPRExists            = errors.New("pull request already exists")
	ErrPRNotFound          = errors.New("pull request not found")
	ErrPRMerged            = errors.New("pull request already merged")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned")
	ErrNoCandidate         = errors.New("no active candidate for reassignment")
)
