package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type shippingaddressseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewShippingAddressSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *shippingaddressseeder {
	return &shippingaddressseeder{db: db, ctx: ctx, logger: logger}
}

func (r *shippingaddressseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM shipping_addresses").Scan(&count)
	if count > 0 {
		r.logger.Debug("shipping_addresses already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO shipping_addresses (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed shipping_addresses", zap.Error(err))
		return err
	}

	r.logger.Info("shipping_addresses seeded successfully")
	return nil
}
