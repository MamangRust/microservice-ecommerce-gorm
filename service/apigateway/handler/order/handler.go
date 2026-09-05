package orderhandler

import (
	
	order_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/order"
		pborder "github.com/MamangRust/microservice-ecommerce-grpc/pb/order"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/order"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
)

type DepsOrder struct {
	Client      *grpc.ClientConn
	StatsClient *grpc.ClientConn
	E           *echo.Echo
	Logger      logger.LoggerInterface
	CacheStore  *cache.CacheStore
}

func RegisterOrderHandler(deps *DepsOrder) {
	mapper := apimapper.NewOrderResponseMapper()
	cache := order_cache.OrderNewMencache(deps.CacheStore)

	queryClient := pborder.NewOrderQueryServiceClient(deps.Client)
	// Stats routes are served by the stats-reader (ClickHouse) service (F3/F5);
	// the legacy OLTP order stats service was removed in F5.
	statsClient := pborder.NewOrderStatsServiceClient(deps.StatsClient)

	NewOrderQueryHandleApi(&orderQueryHandleDeps{
		client: queryClient,
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewOrderCommandHandleApi(&orderCommandHandleDeps{
		client: pborder.NewOrderCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})

	NewOrderStatsHandleApi(&orderStatsHandleDeps{
		client:            statsClient,
		router:            deps.E,
		logger:            deps.Logger,
		mapper:            mapper.StatsMapper(),
		cache:             cache,
		merchantStatsCache: cache,
	})
}
