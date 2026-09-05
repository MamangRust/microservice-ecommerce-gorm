package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type merchantseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *merchantseeder {
	return &merchantseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM merchants").Scan(&count)
	if count > 0 {
		r.logger.Debug("merchants already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO merchants (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed merchants", zap.Error(err))
		return err
	}

	r.logger.Info("merchants seeded successfully")
	return nil
}
