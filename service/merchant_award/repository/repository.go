package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantAwardQuery   MerchantAwardQueryRepository
	MerchantAwardCommand MerchantAwardCommandRepository
	MerchantQuery        MerchantQueryRepository
}

func NewRepositories(DB *gorm.DB, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantAwardQuery:   NewMerchantAwardQueryRepository(DB),
		MerchantAwardCommand: NewMerchantAwardCommandRepository(DB),
		MerchantQuery:        NewMerchantQueryRepository(merchantQuery),
	}
}
