package service

import (
	"context"
	_ "errors"
	"log"
	_ "log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/models"
	"github.com/mink0ff/pr_service/internal/repository"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type PRService struct {
	txManager *repository.TransactionManager
	prRepo    repository.PullRequestRepository
	userRepo  repository.UserRepository
	teamRepo  repository.TeamRepository
}

func NewPRService(txManager *repository.TransactionManager, prRepo repository.PullRequestRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) *PRService {
	return &PRService{
		txManager: txManager,
		prRepo:    prRepo,
		userRepo:  userRepo,
		teamRepo:  teamRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, req *dto.CreatePRRequest) (*dto.CreatePRResponse, error) {
	var resp *dto.CreatePRResponse

	err := s.txManager.Do(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		txUserRepo := s.userRepo.WithTx(tx)
		txPrRepo := s.prRepo.WithTx(tx)

		pr, err := s.createPullRequest(txCtx, req, txPrRepo)
		if err != nil {
			return err
		}

		author, err := s.getAuthorWithTeamLock(txCtx, req.AuthorID, txUserRepo)
		if err != nil {
			return err
		}

		reviewers := s.selectReviewers(txCtx, author.UserID, author.TeamID, txUserRepo)

		if err := s.assignReviewers(txCtx, pr.PullRequestID, reviewers, txPrRepo); err != nil {
			return err
		}

		resp = &dto.CreatePRResponse{
			PR: dto.PullRequestDTO{
				PullRequestID:     pr.PullRequestID,
				PullRequestName:   pr.PullRequestName,
				AuthorID:          pr.AuthorID,
				Status:            dto.PRStatusOpen,
				AssignedReviewers: reviewers,
				CreatedAt:         &pr.CreatedAt,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PRService) MergePR(ctx context.Context, req *dto.MergePRRequest) (*dto.MergePRResponse, error) {
	pr, err := s.prRepo.GetByID(ctx, req.PullRequestID)
	if err != nil || pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRMerged {
		reviewers, _ := s.prRepo.ListReviewers(ctx, req.PullRequestID)
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

	reviewers, _ := s.prRepo.ListReviewers(ctx, req.PullRequestID)
	return &dto.MergePRResponse{
		PR: mapPullRequestToDTO(pr, reviewers),
	}, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error) {
	prID := req.PullRequestID
	oldUserID := req.OldUserID

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

	var candidates []string
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
		ReplacedBy: newReviewer,
	}, nil
}

func mapPullRequestToDTO(pr *models.PullRequest, reviewers []models.User) dto.PullRequestDTO {
	reviewersStr := make([]string, len(reviewers))
	for i, r := range reviewers {
		reviewersStr[i] = r.UserID
	}

	return dto.PullRequestDTO{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            dto.PRStatus(pr.Status),
		AssignedReviewers: reviewersStr,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func (s *PRService) createPullRequest(ctx context.Context, req *dto.CreatePRRequest, prRepo repository.PullRequestRepository) (*models.PullRequest, error) {
	existing, err := prRepo.GetByID(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPRExists
	}

	author, err := s.userRepo.GetByID(ctx, req.AuthorID)
	if err != nil || author == nil {
		return nil, ErrUserNotFound
	}

	pr := models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          models.PROpen,
		CreatedAt:       time.Now(),
	}

	if err := prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *PRService) getAuthorWithTeamLock(ctx context.Context, authorID string, userRepo repository.UserRepository) (*models.User, error) {
	author, err := userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, ErrUserNotFound
	}

	_, err = userRepo.ListActiveByTeam(ctx, author.TeamID)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *PRService) selectReviewers(ctx context.Context, authorID string, teamID uuid.UUID, userRepo repository.UserRepository) []string {
	users, _ := userRepo.ListActiveByTeam(ctx, teamID)

	// фильтруем авторов
	candidates := make([]string, 0, len(users))
	for _, u := range users {
		if u.UserID != authorID {
			candidates = append(candidates, u.UserID)
		}
	}

	n := len(candidates)
	if n == 0 {
		return nil
	}

	maxAssign := 2
	if n < maxAssign {
		maxAssign = n
	}

	selected := make([]string, 0, maxAssign)
	used := make(map[int]struct{})

	for len(selected) < maxAssign {
		idx := rand.Intn(n)
		if _, ok := used[idx]; ok {
			continue
		}
		used[idx] = struct{}{}
		selected = append(selected, candidates[idx])
	}

	return selected
}

func (s *PRService) assignReviewers(ctx context.Context, prID string, reviewers []string, prRepo repository.PullRequestRepository) error {
	for _, r := range reviewers {
		if err := prRepo.AddReviewer(ctx, prID, r); err != nil {
			log.Printf("failed to add reviewer %s: %v", r, err)
			return err
		}
	}
	return nil
}
