package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userSeeder struct {
	db     *gorm.DB
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserSeeder(db *gorm.DB, hash hash.HashPassword, ctx context.Context, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{db: db, hash: hash, ctx: ctx, logger: logger}
}

func (r *userSeeder) Seed() error {
	var count int
	r.db.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	if count > 0 {
		r.logger.Debug("users already seeded, skipping")
		return nil
	}

	now := time.Now()
	for i := 1; i <= 10; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())
		rawPassword := fmt.Sprintf("password%d", i)

		hashedPassword, err := r.hash.HashPassword(rawPassword)
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		err = r.db.WithContext(r.ctx).Exec(
			"INSERT INTO users (firstname, lastname, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			fmt.Sprintf("User%d", i), fmt.Sprintf("Last%d", i), email, hashedPassword, now, now,
		).Error
		if err != nil {
			r.logger.Error("failed to seed user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to seed user %d: %w", i, err)
		}
	}

	r.logger.Info("users seeded successfully")
	return nil
}
