package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type reviewseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewReviewSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *reviewseeder {
	return &reviewseeder{db: db, ctx: ctx, logger: logger}
}

func (r *reviewseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM reviews").Scan(&count)
	if count > 0 {
		r.logger.Debug("reviews already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO reviews (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed reviews", zap.Error(err))
		return err
	}

	r.logger.Info("reviews seeded successfully")
	return nil
}
