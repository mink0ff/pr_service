package migrate

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, migratePath string) error {
	sqlDB, _ := db.DB()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("create driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migratePath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("create migrate instance: %v", err)
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate up failed: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")

	return nil
}
