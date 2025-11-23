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
	"github.com/mink0ff/pr_service/internal/repository/transaction"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type PRServiceImpl struct {
	prRepo      repository.PullRequestRepository
	userRepo    repository.UserRepository
	teamRepo    repository.TeamRepository
	historyRepo repository.ReviewerHistoryRepository
	txManager   *transaction.Manager
}

func NewPRService(prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	historyRepo repository.ReviewerHistoryRepository,
	txManager *transaction.Manager) PRService {
	return &PRServiceImpl{
		prRepo:      prRepo,
		userRepo:    userRepo,
		teamRepo:    teamRepo,
		historyRepo: historyRepo,
		txManager:   txManager,
	}
}

func (s *PRServiceImpl) CreatePR(ctx context.Context, req *dto.CreatePRRequest) (*dto.CreatePRResponse, error) {
	var resp *dto.CreatePRResponse

	err := s.txManager.Do(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		txUserRepo := s.userRepo.WithTx(tx)
		txPrRepo := s.prRepo.WithTx(tx)
		txHistoryRepo := s.historyRepo.WithTx(tx)

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

		if err := s.logReviewerAssignments(txCtx, txHistoryRepo, pr.PullRequestID, reviewers); err != nil {
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

func (s *PRServiceImpl) MergePR(ctx context.Context, req *dto.MergePRRequest) (*dto.MergePRResponse, error) {
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

	mergeTime := time.Now()
	newPr := models.PullRequest{
		PullRequestID:   pr.PullRequestID,
		PullRequestName: pr.PullRequestName,
		AuthorID:        pr.AuthorID,
		Status:          models.PRMerged,
		CreatedAt:       pr.CreatedAt,
		MergedAt:        &mergeTime,
	}

	if err := s.prRepo.Update(ctx, newPr); err != nil {
		return nil, err
	}

	reviewers, _ := s.prRepo.ListReviewers(ctx, req.PullRequestID)
	return &dto.MergePRResponse{
		PR: mapPullRequestToDTO(&newPr, reviewers),
	}, nil
}

func (s *PRServiceImpl) ReassignReviewer(ctx context.Context, req *dto.ReassignReviewerRequest) (*dto.ReassignReviewerResponse, error) {
	var resp *dto.ReassignReviewerResponse

	err := s.txManager.Do(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		txPrRepo := s.prRepo.WithTx(tx)
		txUserRepo := s.userRepo.WithTx(tx)
		txHistoryRepo := s.historyRepo.WithTx(tx)

		pr, err := s.getPRForReassign(txCtx, req.PullRequestID, txPrRepo)
		if err != nil {
			return err
		}

		reviewers, oldReviewer, err := s.getOldReviewer(txCtx, pr.PullRequestID, req.OldUserID, txPrRepo)
		if err != nil {
			return err
		}

		newReviewerID, err := s.pickNewReviewer(txCtx, reviewers, oldReviewer.TeamID, pr.AuthorID, txUserRepo)
		if err != nil {
			return err
		}

		if err := s.updateReviewers(txCtx, pr.PullRequestID, req.OldUserID, newReviewerID, txPrRepo); err != nil {
			return err
		}

		if err := s.logReviewerAssignments(txCtx, txHistoryRepo, pr.PullRequestID, []string{newReviewerID}); err != nil {
			return err
		}

		updatedReviewers, _ := txPrRepo.ListReviewers(txCtx, pr.PullRequestID)

		resp = &dto.ReassignReviewerResponse{
			PR:         mapPullRequestToDTO(pr, updatedReviewers),
			ReplacedBy: newReviewerID,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
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

func (s *PRServiceImpl) createPullRequest(ctx context.Context, req *dto.CreatePRRequest, prRepo repository.PullRequestRepository) (*models.PullRequest, error) {
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

func (s *PRServiceImpl) getAuthorWithTeamLock(ctx context.Context, authorID string, userRepo repository.UserRepository) (*models.User, error) {
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

func (s *PRServiceImpl) selectReviewers(ctx context.Context, authorID string, teamID uuid.UUID, userRepo repository.UserRepository) []string {
	users, _ := userRepo.ListActiveByTeam(ctx, teamID)

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

func (s *PRServiceImpl) assignReviewers(ctx context.Context, prID string, reviewers []string, prRepo repository.PullRequestRepository) error {
	for _, r := range reviewers {
		if err := prRepo.AddReviewer(ctx, prID, r); err != nil {
			log.Printf("failed to add reviewer %s: %v", r, err)
			return err
		}
	}
	return nil
}

func (s *PRServiceImpl) getPRForReassign(ctx context.Context, prID string, prRepo repository.PullRequestRepository) (*models.PullRequest, error) {
	pr, err := prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == models.PRMerged {
		return nil, ErrPRMerged
	}

	return pr, nil
}

func (s *PRServiceImpl) getOldReviewer(
	ctx context.Context,
	prID string,
	oldUserID string,
	prRepo repository.PullRequestRepository,
) ([]models.User, *models.User, error) {

	reviewers, err := prRepo.ListReviewers(ctx, prID)
	if err != nil {
		return nil, nil, err
	}

	var oldReviewer *models.User
	for _, r := range reviewers {
		if r.UserID == oldUserID {
			oldReviewer = &r
			break
		}
	}

	if oldReviewer == nil {
		return nil, nil, ErrReviewerNotAssigned
	}

	return reviewers, oldReviewer, nil
}

func (s *PRServiceImpl) pickNewReviewer(
	ctx context.Context,
	reviewers []models.User,
	teamID uuid.UUID,
	authorID string,
	userRepo repository.UserRepository,
) (string, error) {

	users, err := userRepo.ListActiveByTeam(ctx, teamID)
	if err != nil {
		return "", err
	}

	assigned := map[string]struct{}{}
	for _, r := range reviewers {
		assigned[r.UserID] = struct{}{}
	}

	var candidates []string
	for _, u := range users {

		if u.UserID == authorID {
			continue
		}

		if _, exists := assigned[u.UserID]; exists {
			continue
		}

		candidates = append(candidates, u.UserID)
	}

	if len(candidates) == 0 {
		return "", ErrNoCandidate
	}

	rand.Seed(time.Now().UnixNano())
	return candidates[rand.Intn(len(candidates))], nil
}

func (s *PRServiceImpl) updateReviewers(
	ctx context.Context,
	prID string,
	oldUserID string,
	newUserID string,
	prRepo repository.PullRequestRepository,
) error {

	if err := prRepo.RemoveReviewer(ctx, prID, oldUserID); err != nil {
		return err
	}

	if err := prRepo.AddReviewer(ctx, prID, newUserID); err != nil {
		return err
	}

	return nil
}

func (s *PRServiceImpl) logReviewerAssignments(
	ctx context.Context,
	txRepo repository.ReviewerHistoryRepository,
	prID string,
	reviewers []string,
) error {

	for _, reviewerID := range reviewers {
		event := models.ReviewerAssignmentHistory{
			AssigmentHistoryID: uuid.New(),
			PrID:               prID,
			UserID:             reviewerID,
			CreatedAt:          time.Now(),
		}

		if err := txRepo.AddEvent(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
