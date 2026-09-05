package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type categoryseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *categoryseeder {
	return &categoryseeder{db: db, ctx: ctx, logger: logger}
}

func (r *categoryseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM categories").Scan(&count)
	if count > 0 {
		r.logger.Debug("categories already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO categories (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed categories", zap.Error(err))
		return err
	}

	r.logger.Info("categories seeded successfully")
	return nil
}
