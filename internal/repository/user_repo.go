package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Printf("Failed to create user %v: %v\n", user.UserID, err)
	} else {
		log.Printf("User %v created successfully\n", user.UserID)
	}
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&user, "user_id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("User %v not found\n", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error fetching user %v: %v\n", id, err)
	} else {
		log.Printf("User %v fetched successfully\n", id)
	}

	return &user, err
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ?", user.UserID).
		Save(&user).Error
	if err != nil {
		log.Printf("Failed to update user %v: %v\n", user.UserID, err)
	} else {
		log.Printf("User %v updated successfully\n", user.UserID)
	}
	return err
}

func (r *UserRepo) ListActiveByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("team_id = ? AND is_active = TRUE", teamID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&users).Error
	if err != nil {
		log.Printf("Failed to list active users for team %v: %v\n", teamID, err)
	} else {
		log.Printf("Found %d active users for team %v\n", len(users), teamID)
	}
	return users, err
}

func (r *UserRepo) ListReviewPRs(ctx context.Context, userID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	err := r.db.WithContext(ctx).
		Joins("JOIN pr_reviewers pr ON pr.pull_request_id = pull_requests.pull_request_id").
		Where("pr.reviewer_id = ?", userID).
		Find(&prs).Error
	if err != nil {
		log.Printf("Failed to list PRs for reviewer %v: %v\n", userID, err)
	} else {
		log.Printf("Found %d PRs for reviewer %v\n", len(prs), userID)
	}
	return prs, err
}

func (r *UserRepo) WithTx(tx *gorm.DB) UserRepository {
	log.Println("Creating UserRepository with transaction")
	return &UserRepo{db: tx}
}
