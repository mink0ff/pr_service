package utils

import (
	"github.com/mink0ff/pr_service/internal/repository"
	"github.com/mink0ff/pr_service/internal/repository/transaction"
	"github.com/mink0ff/pr_service/internal/service"
	"gorm.io/gorm"
)

type TestServices struct {
	UserService  service.UserService
	TeamService  service.TeamService
	PRService    service.PRService
	StatsService service.StatsService
	DB           *gorm.DB
	Teardown     func()
}

func InitTestServices() *TestServices {
	db := InitTestDB()

	userRepo := repository.NewUserRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	prRepo := repository.NewPrRepo(db)
	historyRepo := repository.NewReviewerHistoryRepo(db)

	txManager := transaction.NewTransactionManager(db)

	userSvc := service.NewUserService(userRepo, teamRepo)
	teamSvc := service.NewTeamService(teamRepo, userRepo, prRepo, txManager)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo, historyRepo, txManager)
	statsSvc := service.NewStatsService(historyRepo)

	return &TestServices{
		UserService:  userSvc,
		TeamService:  teamSvc,
		PRService:    prSvc,
		StatsService: statsSvc,
		DB:           db,
		Teardown: func() {
			TruncateTables(db)
		},
	}
}
