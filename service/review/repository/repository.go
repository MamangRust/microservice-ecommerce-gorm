package repository

import (
	pbproducts "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbusers "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"gorm.io/gorm"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User    []adapter.GuardOption
	Product []adapter.GuardOption
}

type Repositories struct {
	ProductQuery  ProductQueryRepository
	ReviewQuery   ReviewQueryRepository
	UserQuery     UserQueryRepository
	ReviewCommand ReviewCommandRepository
}

func NewRepositories(db *gorm.DB,
	userQueryClient pbusers.UserQueryServiceClient,
	productQueryClient pbproducts.ProductQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:  productadapter.New(productQueryClient, nil, g.Product...),
		ReviewQuery:   NewReviewQueryRepository(db),
		UserQuery:     useradapter.New(userQueryClient, nil, g.User...),
		ReviewCommand: NewReviewCommandRepository(db),
	}
}
