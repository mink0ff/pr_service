package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/mink0ff/pr_service/internal/models"
	"gorm.io/gorm"
)

type TeamRepo struct {
	db *gorm.DB
}

func NewTeamRepo(db *gorm.DB) TeamRepository {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team models.Team) error {
	err := r.db.WithContext(ctx).Create(&team).Error
	if err != nil {
		log.Printf("Failed to create team %v: %v\n", team.TeamName, err)
	} else {
		log.Printf("Team %v created successfully\n", team.TeamName)
	}
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).First(&team, "team_id = ?", teamID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Team %v not found\n", teamID)
		return nil, nil
	}
	if err != nil {
		log.Printf("Error fetching team %v: %v\n", teamID, err)
	} else {
		log.Printf("Team %v fetched successfully\n", team.TeamName)
	}
	return &team, err
}

func (r *TeamRepo) GetByName(ctx context.Context, teamName string) (*models.Team, error) {
	var team models.Team
	err := r.db.WithContext(ctx).First(&team, "team_name = ?", teamName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Team %v not found\n", teamName)
		return nil, nil
	}
	if err != nil {
		log.Printf("Error fetching team %v: %v\n", teamName, err)
	} else {
		log.Printf("Team %v fetched successfully\n", teamName)
	}
	return &team, err
}

func (r *TeamRepo) ListUsersByTeam(ctx context.Context, teamID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("No users found for team %v\n", teamID)
		return []models.User{}, nil
	}
	if err != nil {
		log.Printf("Failed to list users for team %v: %v\n", teamID, err)
	} else {
		log.Printf("Found %d users for team %v\n", len(users), teamID)
	}
	return users, err
}

func (r *TeamRepo) WithTx(tx *gorm.DB) TeamRepository {
	log.Println("Creating TeamRepository with transaction")
	return &TeamRepo{db: tx}
}
