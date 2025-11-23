package service

import (
	"context"

	"github.com/mink0ff/pr_service/internal/dto"
	"github.com/mink0ff/pr_service/internal/repository"
)

type StatsServiceImpl struct {
	historyRepo repository.ReviewerHistoryRepository
}

func NewStatsService(historyRepo repository.ReviewerHistoryRepository) StatsService {
	return &StatsServiceImpl{historyRepo: historyRepo}
}

func (s *StatsServiceImpl) GetReviewerStats(ctx context.Context) (*dto.ReviewerStatsResponse, error) {
	items, err := s.historyRepo.CountAssignmentsByUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.ReviewerStatsResponse{
		Items: items,
	}, nil
}
