package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type sliderseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewSliderSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *sliderseeder {
	return &sliderseeder{db: db, ctx: ctx, logger: logger}
}

func (r *sliderseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM sliders").Scan(&count)
	if count > 0 {
		r.logger.Debug("sliders already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO sliders (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed sliders", zap.Error(err))
		return err
	}

	r.logger.Info("sliders seeded successfully")
	return nil
}
