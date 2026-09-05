package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type merchantawardseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantAwardSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *merchantawardseeder {
	return &merchantawardseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantawardseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM merchant_awards").Scan(&count)
	if count > 0 {
		r.logger.Debug("merchant_awards already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO merchant_awards (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed merchant_awards", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_awards seeded successfully")
	return nil
}
