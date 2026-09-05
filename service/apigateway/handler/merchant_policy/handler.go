package merchantpolicyhandler

import (
	merchantpolicy_cache "github.com/MamangRust/microservice-ecommerce-grpc-apigateway/cache/merchant_policies"
		pbmerchant_policy "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_policy"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_policy"
	merchantapimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant"
	"github.com/MamangRust/microservice-ecommerce-shared/cache"
	"github.com/MamangRust/microservice-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantPolicy struct {
	Client        *grpc.ClientConn
	E             *echo.Echo
	Logger        logger.LoggerInterface
	CacheStore    *cache.CacheStore
	Observability observability.TraceLoggerObservability
}

func RegisterMerchantPolicyHandler(deps *DepsMerchantPolicy) {
	mapper := apimapper.NewMerchantPolicyResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantpolicy_cache.NewMerchantPoliciesMencache(deps.CacheStore)

	NewMerchantPolicyQueryHandleApi(&merchantPolicyQueryHandleDeps{
		client:        pbmerchant_policy.NewMerchantPolicyQueryServiceClient(deps.Client),
		router:        deps.E,
		logger:        deps.Logger,
		mapper:        mapper.QueryMapper(),
		cache:         cache,
		observability: deps.Observability,
	})

	NewMerchantPolicyCommandHandleApi(&merchantPolicyCommandHandleDeps{
		client:         pbmerchant_policy.NewMerchantPolicyCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
		observability:  deps.Observability,
	})
}
