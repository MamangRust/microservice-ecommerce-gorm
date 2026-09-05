package seeder

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRoleSeeder struct {
	userDB *gorm.DB
	roleDB *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(userDB, roleDB *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{userDB: userDB, roleDB: roleDB, ctx: ctx, logger: logger}
}

func (r *userRoleSeeder) Seed() error {
	var userCount int
	r.userDB.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&userCount)
	var roleCount int
	r.roleDB.WithContext(r.ctx).Raw("SELECT COUNT(*) FROM roles").Scan(&roleCount)

	if userCount == 0 || roleCount == 0 {
		r.logger.Debug("no users or roles found, skipping user_role seeding")
		return nil
	}

	var userIDs []int
	r.userDB.WithContext(r.ctx).Raw("SELECT user_id FROM users WHERE deleted_at IS NULL LIMIT 20").Scan(&userIDs)
	var roleIDs []int
	r.roleDB.WithContext(r.ctx).Raw("SELECT role_id FROM roles LIMIT 4").Scan(&roleIDs)

	now := time.Now()
	for _, uid := range userIDs {
		rid := roleIDs[uid%len(roleIDs)]
		err := r.roleDB.WithContext(r.ctx).Exec(
			"INSERT INTO user_roles (user_id, role_id, created_at, updated_at) VALUES (?, ?, ?, ?) ON CONFLICT DO NOTHING",
			uid, rid, now, now,
		).Error
		if err != nil {
			r.logger.Error("failed to assign role", zap.Error(err))
			return err
		}
	}

	r.logger.Info("user roles seeded successfully")
	return nil
}
