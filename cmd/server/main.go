package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mink0ff/pr_service/internal/handler"
	"github.com/mink0ff/pr_service/internal/repository"
	"github.com/mink0ff/pr_service/internal/repository/gorm"
	"github.com/mink0ff/pr_service/internal/repository/migrate"
	"github.com/mink0ff/pr_service/internal/repository/transaction"
	"github.com/mink0ff/pr_service/internal/service"
)

func main() {
	cfg := &gorm.GormConfig{
		DSN:             "postgres://postgres:postgres@localhost:5432/pr_service?sslmode=disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Minute * 5,
	}

	db, err := gorm.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	migrate.RunMigrations(db)

	userRepo := repository.NewUserRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	prRepo := repository.NewPrRepo(db)

	txManager := transaction.NewTransactionManager(db)

	userService := service.NewUserService(userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo, userRepo, txManager)
	prService := service.NewPRService(prRepo, userRepo, teamRepo, txManager)

	r := chi.NewRouter()
	handler.RegisterRoutes(r, teamService, userService, prService)

	addr := ":8080"
	log.Printf("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
