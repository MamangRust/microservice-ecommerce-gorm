package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	CartQuery    CartQueryRepository
	CartCommand  CartCommandRepository
	UserQuery    UserQueryRepository
	ProductQuery ProductQueryRepository
}

func NewRepositories(DB *gorm.DB,
	userQuery pbuser.UserQueryServiceClient,
	productQuery pbproduct.ProductQueryServiceClient,
) *Repositories {
	return &Repositories{
		CartQuery:    NewCartQueryRepository(DB),
		CartCommand:  NewCartCommandRepository(DB),
		UserQuery:    NewUserQueryRepository(userQuery),
		ProductQuery: NewProductQueryRepository(productQuery),
	}
}
