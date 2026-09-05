package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantBusinessQuery   MerchantBusinessQueryRepository
	MerchantBusinessCommand MerchantBusinessCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(DB *gorm.DB, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantBusinessQuery:   NewMerchantBusinessQueryRepository(DB),
		MerchantBusinessCommand: NewMerchantBusinessCommandRepository(DB),
		MerchantQuery:           NewMerchantQueryRepository(merchantQuery),
	}
}
