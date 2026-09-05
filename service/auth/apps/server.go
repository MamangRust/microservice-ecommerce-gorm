package apps

import (
	pbauth "github.com/MamangRust/microservice-ecommerce-grpc/pb/auth"
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"context"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-auth/cache"
	"github.com/MamangRust/microservice-ecommerce-auth/handler"
	"github.com/MamangRust/microservice-ecommerce-auth/repository"
	"github.com/MamangRust/microservice-ecommerce-auth/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/auth"
	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pkgresilience "github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"google.golang.org/grpc/credentials/insecure"
)

// kafkaOutboxPublisher adapts the ecommerce *kafka.Kafka (whose SendMessage
// takes no context) to the outbox.OutboxPublisher contract.
type kafkaOutboxPublisher struct {
	k *kafka.Kafka
}

func (p kafkaOutboxPublisher) SendMessage(_ context.Context, topic, key string, value []byte) error {
	return p.k.SendMessage(topic, key, value)
}

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")
	userAddr := viper.GetString("GRPC_USER_ADDR")

	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(pkgresilience.NewDependencyGuardInterceptor(srv.Logger).UnaryInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(pkgresilience.NewDependencyGuardInterceptor(srv.Logger).UnaryInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	roleQueryClient := pbrole.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pbrole.NewRoleCommandServiceClient(roleConn)
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)
	userCommandClient := pbuser.NewUserCommandServiceClient(userConn)

	hasher := hash.NewHashingPassword()
	repositories := repository.NewRepositories(srv.GormDB, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	obs, _ := observability.NewObservability("auth-server", srv.Logger)

	c := cache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(service.NewGormOutboxQuerier(srv.GormDB), kafkaOutboxPublisher{k: myKafka}, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      c,
		Repositories:  repositories,
		Token:         tokenManager,
		Hash:          hasher,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		GormDB:        srv.GormDB,
		Outbox:        outboxService,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{Service: services, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbauth.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
