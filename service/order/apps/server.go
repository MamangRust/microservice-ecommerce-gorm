package apps

import (
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc/pb/order"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-order/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-order/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-order/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-order/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pkgresilience "github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// gRPC Client Connections. F6: dependency guard (per-call deadline + circuit
	// breaker + bulkhead) on every downstream gRPC dependency.
	guard := pkgresilience.NewDependencyGuardInterceptor(srv.Logger)

	userAddr := viper.GetString("GRPC_USER_ADDR")

	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)
	productCommandClient := pbproduct.NewProductCommandServiceClient(productConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "merchant:50055"
	}
	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	orderItemAddr := viper.GetString("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "order-item:50056"
	}
	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order_item service: %w", err)
	}
	orderItemQueryClient := pborder_item.NewOrderItemQueryServiceClient(orderItemConn)
	orderItemCommandClient := pborder_item.NewOrderItemCommandServiceClient(orderItemConn)

	shippingAddr := viper.GetString("GRPC_SHIPPING_ADDRESS_ADDR")
	if shippingAddr == "" {
		shippingAddr = "shipping_address:50063"
	}
	shippingConn, err := grpc.NewClient(shippingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping_address service: %w", err)
	}
	shippingCommandClient := pbshipping_address.NewShippingCommandServiceClient(shippingConn)

	transactionAddr := viper.GetString("GRPC_TRANSACTION_ADDR")
	if transactionAddr == "" {
		transactionAddr = "transaction:50061"
	}
	transactionConn, err := grpc.NewClient(transactionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to transaction service: %w", err)
	}
	transactionCommandClient := pbtransaction.NewTransactionCommandServiceClient(transactionConn)

	repos := repository.NewRepositories(&repository.Deps{
		DB:                 srv.GormDB,
		UserQuery:          userQueryClient,
		ProductQuery:       productQueryClient,
		ProductCommand:     productCommandClient,
		MerchantQuery:      merchantQueryClient,
		OrderItemQuery:     orderItemQueryClient,
		OrderItemCommand:   orderItemCommandClient,
		ShippingCommand:    shippingCommandClient,
		TransactionCommand: transactionCommandClient,
	})

	obs, _ := observability.NewObservability("order-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	svc := service.NewService(&service.Deps{
		Kafka:         myKafka,
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	// Start the outbox relay so stats events committed with the order flow are
	// published to Kafka with durable retry and dead-letter semantics (F3).
	go svc.Outbox.Start(srv.Ctx, service.OutboxRelayInterval, service.OutboxRelayBatchSize)

	srv.RegisterServices = func(gs *grpc.Server) {
		pborder.RegisterOrderQueryServiceServer(gs, h.OrderQuery)
		pborder.RegisterOrderCommandServiceServer(gs, h.OrderCommand)
	}

	return srv, nil
}
