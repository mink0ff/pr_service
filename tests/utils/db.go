package utils

import (
	"log"

	"github.com/mink0ff/pr_service/internal/config"
	"github.com/mink0ff/pr_service/internal/repository/gormdb"
	"github.com/mink0ff/pr_service/internal/repository/migrate"
	"gorm.io/gorm"
)

// InitTestDB подключается к тестовой базе и прогоняет миграции
func InitTestDB() *gorm.DB {
	cfg := config.LoadDBConfig("../../.env.test")

	db, err := gormdb.NewGormDB(&gormdb.GormConfig{
		DSN:             cfg.DSN,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLife,
	})
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	if err := migrate.RunMigrations(db, "../../migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

// TruncateTables очищает все тестовые таблицы
func TruncateTables(db *gorm.DB) {
	tables := []string{"users", "teams", "pull_requests", "pr_reviewers", "reviewer_assignment_histories"}

	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE" + " " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}
