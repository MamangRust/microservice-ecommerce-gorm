package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
)

type productQueryRepository struct {
	client pbproduct.ProductQueryServiceClient
}

func NewProductQueryRepository(client pbproduct.ProductQueryServiceClient) ProductQueryRepository {
	return &productQueryRepository{client: client}
}

func (r *productQueryRepository) FindById(ctx context.Context, id int) (*ProductResult, error) {
	res, err := r.client.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}

	return &ProductResult{
		ProductID:    res.Data.Id,
		Name:         res.Data.Name,
		Price:        res.Data.Price,
		CountInStock: res.Data.CountInStock,
		ImageProduct: &res.Data.ImageProduct,
		Weight:       &res.Data.Weight,
	}, nil
}
