package apps

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"fmt"

	"github.com/MamangRust/microservice-ecommerce-grpc-user/cache"
	"github.com/MamangRust/microservice-ecommerce-grpc-user/handler"
	"github.com/MamangRust/microservice-ecommerce-grpc-user/repository"
	"github.com/MamangRust/microservice-ecommerce-grpc-user/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/hash"
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

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")

	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(pkgresilience.NewDependencyGuardInterceptor(srv.Logger).UnaryInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	roleClient := pbrole.NewRoleQueryServiceClient(roleConn)

	repos := repository.NewRepositories(srv.GormDB, roleClient)
	hashing := hash.NewHashingPassword()
	obs, _ := observability.NewObservability("user-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Repositories:  repos,
		Hash:          hashing,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbuser.RegisterUserQueryServiceServer(gs, h.UserQuery)
		pbuser.RegisterUserCommandServiceServer(gs, h.UserCommand)
	}

	return srv, nil
}
