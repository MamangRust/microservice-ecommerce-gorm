package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type roleseeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *roleseeder {
	return &roleseeder{db: db, ctx: ctx, logger: logger}
}

func (r *roleseeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM roles").Scan(&count)
	if count > 0 {
		r.logger.Debug("roles already seeded, skipping")
		return nil
	}

	now := time.Now()
	err := r.db.WithContext(r.ctx).Exec(
		"INSERT INTO roles (name, slug_category, created_at, updated_at) VALUES ('Electronics', 'electronics', ?, ?)",
		now, now,
	).Error
	if err != nil {
		r.logger.Error("failed to seed roles", zap.Error(err))
		return err
	}

	r.logger.Info("roles seeded successfully")
	return nil
}
