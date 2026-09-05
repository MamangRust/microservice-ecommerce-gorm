package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryRepository
	MerchantPoliciesCommand MerchantPoliciesCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(DB *gorm.DB, merchantQuery pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantPoliciesQuery:   NewMerchantPolicyQueryRepository(DB),
		MerchantPoliciesCommand: NewMerchantPolicyCommandRepository(DB),
		MerchantQuery:           NewMerchantQueryRepository(merchantQuery),
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
