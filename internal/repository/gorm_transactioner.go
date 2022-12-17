package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	// GormTransactioner :nodoc:
	GormTransactioner interface {
		Begin(ctx context.Context) *gorm.DB
		Commit(tx *gorm.DB) error
		Rollback(tx *gorm.DB)
	}

	gormTransactioner struct {
		db *gorm.DB
	}
)

// NewGormTransactioner :nodoc:
func NewGormTransactioner(db *gorm.DB) GormTransactioner {
	return &gormTransactioner{db: db}
}

// Begin :nodoc:
func (t *gormTransactioner) Begin(ctx context.Context) *gorm.DB {
	return t.db.WithContext(ctx).Begin()
}

// Commit :nodoc:
func (t *gormTransactioner) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// Rollback :nodoc:
func (t *gormTransactioner) Rollback(tx *gorm.DB) {
	tx.Rollback()
}
