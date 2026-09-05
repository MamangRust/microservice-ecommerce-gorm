package transactionhandler

import (
	
	transaction_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/transaction"
		pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/transaction"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
)

type DepsTransaction struct {
	Client      *grpc.ClientConn
	StatsClient *grpc.ClientConn
	E           *echo.Echo
	Logger      logger.LoggerInterface
	CacheStore  *cache.CacheStore
	ApiHandler  sharedErrors.ApiHandler
}

func RegisterTransactionHandler(deps *DepsTransaction) {
	mapper := apimapper.NewTransactionResponseMapper()
	statsMapper := apimapper.NewTransactionStatsResponseMapper()
	cache := transaction_cache.NewTransactionMencache(deps.CacheStore)

	queryClient := pbtransaction.NewTransactionQueryServiceClient(deps.Client)
	commandClient := pbtransaction.NewTransactionCommandServiceClient(deps.Client)
	// Stats routes are served by the stats-reader (ClickHouse) service (F3/F5);
	// the legacy OLTP transaction stats service was removed in F5.
	statsClient := pbtransaction.NewTransactionStatsServiceClient(deps.StatsClient)
	statsByMerchantClient := pbtransaction.NewTransactionStatsByMerchantServiceClient(deps.StatsClient)

	NewTransactionQueryHandleApi(&transactionQueryHandleDeps{
		queryClient: queryClient,
		router:      deps.E,
		logger:      deps.Logger,
		mapper:      mapper.QueryMapper(),
		cache:       cache,
		apiHandler:  deps.ApiHandler,
	})

	NewTransactionCommandHandleApi(&transactionCommandHandleDeps{
		client:     commandClient,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewTransactionStatsHandleApi(&transactionStatsHandleDeps{
		statsClient:           statsClient,
		statsByMerchantClient: statsByMerchantClient,
		router:                deps.E,
		logger:                deps.Logger,
		statsMapper:           statsMapper,
		statsCache:            cache,
		statsByMerchantCache:  cache,
		apiHandler:            deps.ApiHandler,
	})
}
