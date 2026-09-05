package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

func NewRepositories(DB *gorm.DB, userQueryClient pbuser.UserQueryServiceClient, productQueryClient pbproduct.ProductQueryServiceClient) *Repositories {
	return &Repositories{
		ProductQuery:  NewProductQueryRepository(productQueryClient),
		ReviewQuery:   NewReviewQueryRepository(DB),
		UserQuery:     NewUserQueryRepository(userQueryClient),
		ReviewCommand: NewReviewCommandRepository(DB),
	}
}
