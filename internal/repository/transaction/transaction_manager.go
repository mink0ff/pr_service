package transaction

import (
	"context"

	"gorm.io/gorm"
)

type Manager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *Manager {
	return &Manager{db: db}
}

func (t *Manager) Do(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}
