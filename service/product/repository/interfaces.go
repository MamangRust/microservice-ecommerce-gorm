package repository

import (
	"context"
	"time"

	categoryadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// CategoryQueryRepository is the shared category query contract, provided by
// the category gRPC adapter.
type CategoryQueryRepository = categoryadapter.QueryRepository

// MerchantQueryRepository is the shared merchant query contract, provided by the
// merchant gRPC adapter so this service never holds a raw gRPC client.
type MerchantQueryRepository = merchantadapter.QueryRepository

type ProductResult struct {
	ProductID    int32
	MerchantID   int32
	CategoryID   int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	Weight       *int32
	Rating       *float64
	SlugProduct  *string
	ImageProduct *string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type ProductQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*ProductResult, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*ProductResult, error)
	FindByID(ctx context.Context, product_id int) (*ProductResult, error)
	FindByIDTrashed(ctx context.Context, product_id int) (*ProductResult, error)
}

type ProductCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateProductRequest) (*models.Product, error)
	Update(ctx context.Context, request *requests.UpdateProductRequest) (*models.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*models.Product, error)
	AdjustProductStock(ctx context.Context, product_id int, delta int, operationID string) (*models.Product, error)
	CleanupStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error)
	Trash(ctx context.Context, product_id int) (*models.Product, error)
	Restore(ctx context.Context, product_id int) (*models.Product, error)
	DeletePermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
