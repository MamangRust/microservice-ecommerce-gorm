package repository

import (
	pbcategories "github.com/MamangRust/microservice-ecommerce-grpc-pb/category"
	pbmerchants "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	categoryadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"gorm.io/gorm"
)

type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(db *gorm.DB,
	categoryQueryClient pbcategories.CategoryQueryServiceClient,
	merchantQueryClient pbmerchants.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	categoryQuery := categoryadapter.New(categoryQueryClient, g.Category...)

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db, categoryQuery),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  categoryQuery,
		MerchantQuery:  merchantadapter.New(merchantQueryClient, g.Merchant...),
	}
}
