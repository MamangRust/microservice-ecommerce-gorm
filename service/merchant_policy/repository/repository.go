package repository

import (
	pbmerchants "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"gorm.io/gorm"
)

type GuardOptions struct {
	Merchant []adapter.GuardOption
}

type Repositories struct {
	MerchantPoliciesQuery   MerchantPoliciesQueryRepository
	MerchantPoliciesCommand MerchantPoliciesCommandRepository
	MerchantQuery           MerchantQueryRepository
}

func NewRepositories(db *gorm.DB, merchantQueryClient pbmerchants.MerchantQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantPoliciesQuery:   NewMerchantPolicyQueryRepository(db),
		MerchantPoliciesCommand: NewMerchantPolicyCommandRepository(db),
		MerchantQuery:           merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
