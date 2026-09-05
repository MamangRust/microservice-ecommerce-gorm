package main

import (
	"gorm.io/gorm"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/seeder"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

// open connects to the database configured via the given DBCluster prefix.
func open(logger logger.LoggerInterface, prefix string) (*gorm.DB, func(), error) {
	gormDB, err := database.NewGormClientWithPrefix(logger, prefix)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to %s: %w", prefix, err)
	}
	closeFn := func() {
		if sqlDB, closeErr := gormDB.DB(); closeErr == nil {
			sqlDB.Close()
		}
	}
	return gormDB, closeFn, nil
}

// openUser connects to the user database, which the role seeder reads across
// the per-service database boundary (F-per-service schema).
func openUser(logger logger.LoggerInterface, prefix string) (*gorm.DB, func(), error) {
	gormDB, err := database.NewGormClientWithPrefix(logger, prefix)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to %s: %w", prefix, err)
	}
	closeFn := func() {
		if sqlDB, closeErr := gormDB.DB(); closeErr == nil {
			sqlDB.Close()
		}
	}
	return gormDB, closeFn, nil
}

func main() {
	logger, err := logger.NewLogger("seeder", nil)
	if err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	ctx := context.Background()

	roleDB, closeRole, err := open(logger, "DB_ROLE")
	if err != nil {
		logger.Fatal("Failed to connect to role database", zap.Error(err))
	}
	defer closeRole()

	userDB, closeUser, err := openUser(logger, "DB_USER")
	if err != nil {
		logger.Fatal("Failed to connect to user database", zap.Error(err))
	}
	defer closeUser()

	role := seeder.NewRoleSeeder(roleDB, ctx, logger)
	if err := role.Seed(); err != nil {
		logger.Fatal("Failed to seed roles", zap.Error(err))
	}

	userRole := seeder.NewUserRoleSeeder(userDB, roleDB, ctx, logger)
	if err := userRole.Seed(); err != nil {
		logger.Fatal("Failed to seed user_roles", zap.Error(err))
	}

	logger.Info("roles and user_roles seeded successfully")
}
