package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type merchantbusinessseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantBusinessSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *merchantbusinessseeder {
	return &merchantbusinessseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantbusinessseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM merchant_businesses").Scan(&count)
	if count > 0 {
		r.logger.Debug("merchant_businesses already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO merchant_businesses (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed merchant_businesses", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_businesses seeded successfully")
	return nil
}
