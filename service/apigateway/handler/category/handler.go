package categoryhandler

import (
	
	category_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/category"
		pbcategory "github.com/MamangRust/microservice-ecommerce-grpc/pb/category"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-pkg/upload_image"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/category"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsCategory struct {
	Client      *grpc.ClientConn
	StatsClient *grpc.ClientConn
	E           *echo.Echo
	Logger      logger.LoggerInterface
	CacheStore  *cache.CacheStore
	UploadImage upload_image.ImageUploads
	ApiHandler  errors.ApiHandler
}

func RegisterCategoryHandler(deps *DepsCategory) {
	mapper := apimapper.NewCategoryResponseMapper()
	cache := category_cache.NewCategoryMencache(deps.CacheStore)

	handlers := []func(){
		setupCategoryQueryHandler(deps, mapper.QueryMapper(), cache),
		setupCategoryCommandHandler(deps, mapper.CommandMapper(), cache),
		setupCategoryStatsHandler(deps, mapper.StatsMapper(), cache),
	}

	for _, h := range handlers {
		h()
	}
}

func setupCategoryQueryHandler(deps *DepsCategory, mapper apimapper.CategoryQueryResponseMapper, cache category_cache.CategoryMencache) func() {
	return func() {
		NewCategoryQueryHandleApi(&categoryQueryHandleDeps{
			client:     pbcategory.NewCategoryQueryServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupCategoryCommandHandler(deps *DepsCategory, mapper apimapper.CategoryCommandResponseMapper, cache category_cache.CategoryMencache) func() {
	return func() {
		NewCategoryCommandHandleApi(&categoryCommandHandleDeps{
			client:     pbcategory.NewCategoryCommandServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			upload_image: deps.UploadImage,
			apiHandler: deps.ApiHandler,
		})
	}
}
func setupCategoryStatsHandler(deps *DepsCategory, mapper apimapper.CategoryStatsResponseMapper, cache category_cache.CategoryMencache) func() {
	return func() {
		// Stats routes are served by the stats-reader (ClickHouse) service
		// (F3/F5); the legacy OLTP category stats service was removed in F5.
		statsClient := pbcategory.NewCategoryStatsServiceClient(deps.StatsClient)
		statsByIdClient := pbcategory.NewCategoryStatsByIdServiceClient(deps.StatsClient)
		statsByMerchantClient := pbcategory.NewCategoryStatsByMerchantServiceClient(deps.StatsClient)
		NewCategoryStatsHandleApi(&categoryStatsHandleDeps{
			statsClient:           statsClient,
			statsByIdClient:       statsByIdClient,
			statsByMerchantClient: statsByMerchantClient,
			router:                deps.E,
			logger:                deps.Logger,
			mapper:                mapper,
			cache:                 cache,
			apiHandler:            deps.ApiHandler,
		})
	}
}
