package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type bannerseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewBannerSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *bannerseeder {
	return &bannerseeder{db: db, ctx: ctx, logger: logger}
}

func (r *bannerseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM banners").Scan(&count)
	if count > 0 {
		r.logger.Debug("banners already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO banners (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed banners", zap.Error(err))
		return err
	}

	r.logger.Info("banners seeded successfully")
	return nil
}
