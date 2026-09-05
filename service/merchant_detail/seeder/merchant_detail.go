package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type merchantdetailseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantDetailSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *merchantdetailseeder {
	return &merchantdetailseeder{db: db, ctx: ctx, logger: logger}
}

func (r *merchantdetailseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM merchant_details").Scan(&count)
	if count > 0 {
		r.logger.Debug("merchant_details already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO merchant_details (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed merchant_details", zap.Error(err))
		return err
	}

	r.logger.Info("merchant_details seeded successfully")
	return nil
}
