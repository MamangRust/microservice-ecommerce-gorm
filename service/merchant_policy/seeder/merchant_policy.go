package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type merchantpolicyseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantPolicySeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *merchantpolicyseeder {
	return &merchantpolicyseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantpolicyseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM merchant_policies").Scan(&count)
	if count > 0 {
		r.logger.Debug("merchant_policies already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO merchant_policies (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed merchant_policies", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_policies seeded successfully")
	return nil
}
