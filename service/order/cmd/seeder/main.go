package main

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/seeder"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/dotenv"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

func main() {
	l, err := logger.NewLogger("seeder", nil)
	if err != nil {
		panic(err)
	}

	if err := dotenv.Viper(); err != nil {
		l.Fatal("Failed to load .env", zap.Error(err))
	}

	gormDB, err := database.NewGormClientWithPrefix(l, "DB_ORDER")
	if err != nil {
		l.Fatal("Failed to connect", zap.Error(err))
	}

	ctx := context.Background()

	s := seeder.NewOrderSeeder(gormDB, ctx, l)
	if err := s.Seed(); err != nil {
		l.Fatal("Failed to seed", zap.Error(err))
	}

	l.Info("seeded successfully")
}
