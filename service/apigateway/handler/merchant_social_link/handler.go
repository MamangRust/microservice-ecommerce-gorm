package merchantsociallinkhandler

import (
		pbmerchant_social_link "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant_social_link"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/microservice-ecommerce-shared/mapper/merchant_social_link"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantSocialLink struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger     logger.LoggerInterface
	ApiHandler sharedErrors.ApiHandler
}

func RegisterMerchantSocialLinkHandler(deps *DepsMerchantSocialLink) {
	mapper := apimapper.NewMerchantSocialLinkResponseMapper()

	NewMerchantSocialLinkCommandHandleApi(&merchantSocialLinkCommandHandleDeps{
		client: pbmerchant_social_link.NewMerchantSocialCommandServiceClient(deps.Client),
		router: deps.E,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		apiHandler: deps.ApiHandler,
	})
}
