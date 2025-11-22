package service

import (
	"context"
	_ "errors"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
)

type PRService struct {
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(prRepo repository.PullRequestRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, req *dto.CreatePRRequest) (*dto.CreatePRResponse, error) {
	existing, err := s.prRepo.GetByID(ctx, uuid.MustParse(req.PullRequestID))
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPRExists
	}

	authorID := uuid.MustParse(req.AuthorID)
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, ErrUserNotFound
	}

	pr := models.PullRequest{
		PullRequestID:   uuid.MustParse(req.PullRequestID),
		PullRequestName: req.PullRequestName,
		AuthorID:        authorID,
		Status:          models.PROpen,
		CreatedAt:       time.Now(),
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	users, err := s.userRepo.ListActiveByTeam(ctx, author.TeamID)
	if err != nil {
		return nil, err
	}

	var reviewers []uuid.UUID
	for _, u := range users {
		if u.UserID != authorID {
			reviewers = append(reviewers, u.UserID)
		}
	}

	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })

	maxAssign := 2
	if len(reviewers) < 2 {
		maxAssign = len(reviewers)
	}

	assigned := reviewers[:maxAssign]

	for _, r := range assigned {
		if err := s.prRepo.AddReviewer(ctx, pr.PullRequestID, r); err != nil {
			log.Printf("add reviewer error: %v", err)
		}
	}

	assignedStr := make([]string, len(assigned))
	for i, r := range assigned {
		assignedStr[i] = r.String()
	}

	resp := &dto.CreatePRResponse{
		PR: dto.PullRequestDTO{
			PullRequestID:     pr.PullRequestID.String(),
			PullRequestName:   pr.PullRequestName,
			AuthorID:          pr.AuthorID.String(),
			Status:            dto.PRStatusOpen,
			AssignedReviewers: assignedStr,
			CreatedAt:         &pr.CreatedAt,
		},
	}
	return resp, nil
}

func (s *PRService) MergePR(ctx context.Context, req *dto.MergePRRequest) (*dto.MergePRResponse, error) {
	prID := uuid.MustParse(req.PullRequestID)
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRMerged {
		reviewers, _ := s.prRepo.ListReviewers(ctx, prID)
		return &dto.MergePRResponse{
			PR: mapPullRequestToDTO(pr, reviewers),
		}, nil
	}

	now := time.Now()
	pr.Status = models.PRMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		return nil, err
	}

	reviewers, _ := s.prRepo.ListReviewers(ctx, prID)
	return &dto.MergePRResponse{
		PR: mapPullRequestToDTO(pr, reviewers),
	}, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error) {
	prID := uuid.MustParse(req.PullRequestID)
	oldUserID := uuid.MustParse(req.OldUserID)

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRMerged {
		return nil, ErrPRMerged
	}

	reviewers, err := s.prRepo.ListReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	var oldReviewer *models.User
	for _, r := range reviewers {
		if r.UserID == oldUserID {
			oldReviewer = &r
			break
		}
	}
	if oldReviewer == nil {
		return nil, ErrReviewerNotAssigned
	}

	users, err := s.userRepo.ListActiveByTeam(ctx, oldReviewer.TeamID)
	if err != nil {
		return nil, err
	}

	candidates := []uuid.UUID{}
	for _, u := range users {
		skip := false
		for _, r := range reviewers {
			if r.UserID == u.UserID {
				skip = true
				break
			}
		}
		if !skip {
			candidates = append(candidates, u.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoCandidate
	}

	rand.Seed(time.Now().UnixNano())
	newReviewer := candidates[rand.Intn(len(candidates))]

	if err := s.prRepo.RemoveReviewer(ctx, prID, oldUserID); err != nil {
		return nil, err
	}
	if err := s.prRepo.AddReviewer(ctx, prID, newReviewer); err != nil {
		return nil, err
	}

	updatedReviewers, _ := s.prRepo.ListReviewers(ctx, prID)

	return &dto.ReassignReviewerResponse{
		PR:         mapPullRequestToDTO(pr, updatedReviewers),
		ReplacedBy: newReviewer.String(),
	}, nil
}

func mapPullRequestToDTO(pr *models.PullRequest, reviewers []models.User) dto.PullRequestDTO {
	reviewersStr := make([]string, len(reviewers))
	for i, r := range reviewers {
		reviewersStr[i] = r.UserID.String()
	}

	return dto.PullRequestDTO{
		PullRequestID:     pr.PullRequestID.String(),
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID.String(),
		Status:            dto.PRStatus(pr.Status),
		AssignedReviewers: reviewersStr,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
