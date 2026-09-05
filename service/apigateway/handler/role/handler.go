package rolehandler

import (
	api_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/role"
	"github.com/MamangRust/microservice-ecommerce-pkg/kafka"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/role"
		pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsRole struct {
	Kafka      *kafka.Kafka
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	Cache      api_cache.RoleCache
	ApiHandler errors.ApiHandler
}

func RegisterRoleHandler(deps *DepsRole) {
	mapper := apimapper.NewRoleResponseMapper()
	cache := role_cache.NewRoleMencache(deps.CacheStore)

	NewRoleQueryHandleApi(&roleQueryHandleDeps{
		client:     pbrole.NewRoleQueryServiceClient(deps.Client),
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.QueryMapper(),
		kafka:      deps.Kafka,
		cache_role: deps.Cache,
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewRoleCommandHandleApi(&roleCommandHandleDeps{
		client:      pbrole.NewRoleCommandServiceClient(deps.Client),
		queryClient: pbrole.NewRoleQueryServiceClient(deps.Client),
		router:      deps.E,
		logger:      deps.Logger,
		mapper:      mapper.CommandMapper(),
		kafka:       deps.Kafka,
		cache_role:  deps.Cache,
		cache:       cache,
		apiHandler:  deps.ApiHandler,
	})
}
