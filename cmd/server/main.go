package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mink0ff/pr_service/internal/config"
	"github.com/mink0ff/pr_service/internal/handler"
	"github.com/mink0ff/pr_service/internal/repository"
	"github.com/mink0ff/pr_service/internal/repository/gormdb"
	"github.com/mink0ff/pr_service/internal/repository/migrate"
	"github.com/mink0ff/pr_service/internal/repository/transaction"
	"github.com/mink0ff/pr_service/internal/service"
)

func main() {
	cfg := config.LoadDBConfig(".env")

	db, err := gormdb.NewGormDB(&gormdb.GormConfig{
		DSN:             cfg.DSN,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLife,
	})

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err = migrate.RunMigrations(db, cfg.MigrationPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	prRepo := repository.NewPrRepo(db)
	reviewerHistoryPero := repository.NewReviewerHistoryRepo(db)

	txManager := transaction.NewTransactionManager(db)

	userService := service.NewUserService(userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo, userRepo, prRepo, txManager)
	prService := service.NewPRService(prRepo, userRepo, teamRepo, reviewerHistoryPero, txManager)
	statsService := service.NewStatsService(reviewerHistoryPero)

	r := chi.NewRouter()
	handler.RegisterRoutes(r, teamService, userService, prService, statsService)

	addr := ":8080"
	log.Printf("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
