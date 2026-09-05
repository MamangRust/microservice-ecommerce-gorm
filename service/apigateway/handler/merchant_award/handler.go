package merchantawardhandler

import (
	merchantaward_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/merchant_awards"
		pbmerchant_award "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_award"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_award"
	merchantapimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
)

type DepsMerchantAward struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterMerchantAwardHandler(deps *DepsMerchantAward) {
	mapper := apimapper.NewMerchantAwardResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantaward_cache.NewMerchantAward(deps.CacheStore)

	NewMerchantAwardQueryHandleApi(&merchantAwardQueryHandleDeps{
		client: pbmerchant_award.NewMerchantAwardQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantAwardCommandHandleApi(&merchantAwardCommandHandleDeps{
		client:         pbmerchant_award.NewMerchantAwardCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
