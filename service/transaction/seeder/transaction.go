package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type transactionseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *transactionseeder {
	return &transactionseeder{db: db, ctx: ctx, logger: logger}
}

func (r *transactionseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM transactions").Scan(&count)
	if count > 0 {
		r.logger.Debug("transactions already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO transactions (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed transactions", zap.Error(err))
		return err
	}

	r.logger.Info("transactions seeded successfully")
	return nil
}
