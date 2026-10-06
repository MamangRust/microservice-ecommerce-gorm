package main

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	"github.com/MamangRust/microservice-ecommerce-grpc-merchant_policy/seeder"
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

func main() {
	logger, err := logger.NewLogger("seeder", nil)
	if err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	ctx := context.Background()

	q, closeFn, err := open(logger, database.MerchantCluster)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer closeFn()

	s := seeder.NewMerchantPolicySeeder(q, ctx, logger)
	if err := s.Seed(); err != nil {
		logger.Fatal("Failed to seed merchant policies", zap.Error(err))
	}

	logger.Info("merchant policies seeded successfully")
}
