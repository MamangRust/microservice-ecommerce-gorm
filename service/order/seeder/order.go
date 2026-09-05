package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type orderseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *orderseeder {
	return &orderseeder{db: db, ctx: ctx, logger: logger}
}

func (r *orderseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM orders").Scan(&count)
	if count > 0 {
		r.logger.Debug("orders already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO orders (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed orders", zap.Error(err))
		return err
	}

	r.logger.Info("orders seeded successfully")
	return nil
}
