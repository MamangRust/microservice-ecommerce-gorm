package service

import (
	mencache "github.com/MamangRust/microservice-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"gorm.io/gorm"
)

type Service struct {
	TransactionQuery   TransactionQueryService
	TransactionCommand TransactionCommandService
	Outbox             OutboxService
}

type Deps struct {
	Kafka         *kafka.Kafka
	DB            *gorm.DB
	Cache         *mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		TransactionQuery: NewTransactionQueryService(&TransactionQueryServiceDeps{
			Observability: deps.Observability,
			Cache:         deps.Cache.TransactionQueryCache,
			Repository:    deps.Repositories.TransactionQuery,
			Logger:        deps.Logger,
		}),
		TransactionCommand: NewTransactionCommandService(&TransactionCommandServiceDeps{
			Observability:      deps.Observability,
			Kafka:              deps.Kafka,
			DB:                 deps.DB,
			Cache:              deps.Cache.TransactionCommandCache,
			TransactionQuery:   deps.Repositories.TransactionQuery,
			TransactionCommand: deps.Repositories.TransactionCommand,
			UserQuery:          deps.Repositories.UserQuery,
			MerchantQuery:      deps.Repositories.MerchantQuery,
			OrderQuery:         deps.Repositories.OrderQuery,
			OrderItem:          deps.Repositories.OrderItem,
			ShippingAddress:    deps.Repositories.ShippingAddress,
			Outbox:             deps.Repositories.Outbox,
			Logger:             deps.Logger,
		}),
		Outbox: NewOutboxService(deps.Repositories.Outbox, deps.Kafka, deps.Logger),
	}
}
