package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type productseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *productseeder {
	return &productseeder{db: db, ctx: ctx, logger: logger}
}

func (r *productseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		r.logger.Debug("products already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO products (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed products", zap.Error(err))
		return err
	}

	r.logger.Info("products seeded successfully")
	return nil
}
