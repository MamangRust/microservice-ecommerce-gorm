package apps

import (
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-cart/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/server"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
		pbcart "github.com/MamangRust/microservice-ecommerce-grpc/pb/cart"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := viper.GetString("GRPC_USER_ADDR")
	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	productConn, err := grpc.NewClient(
		productAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	userQueryClient := pbuser.NewUserQueryServiceClient(userConn)
	productQueryClient := pbproduct.NewProductQueryServiceClient(productConn)

	repos := repository.NewRepositories(srv.GormDB, userQueryClient, productQueryClient)

	obs, _ := observability.NewObservability("cart-service", srv.Logger)
	mencache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         mencache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbcart.RegisterCartQueryServiceServer(gs, h.CartQuery)
		pbcart.RegisterCartCommandServiceServer(gs, h.CartCommand)
	}

	return srv, nil
}
