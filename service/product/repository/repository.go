package repository

import (
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc/pb/category"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	"gorm.io/gorm"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

type categoryQueryRepository struct {
	client pbcategory.CategoryQueryServiceClient
}

func NewCategoryQueryRepository(client pbcategory.CategoryQueryServiceClient) *categoryQueryRepository {
	return &categoryQueryRepository{client: client}
}

type merchantQueryRepository struct {
	client pbmerchant.MerchantQueryServiceClient
}

func NewMerchantQueryRepository(client pbmerchant.MerchantQueryServiceClient) *merchantQueryRepository {
	return &merchantQueryRepository{client: client}
}

func NewRepositories(db *gorm.DB, categoryClient pbcategory.CategoryQueryServiceClient, merchantClient pbmerchant.MerchantQueryServiceClient) *Repositories {
	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db, categoryClient),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  NewCategoryQueryRepository(categoryClient),
		MerchantQuery:  NewMerchantQueryRepository(merchantClient),
	}
}
