package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-review/dto"
	product_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
)

type productQueryRepository struct {
	client pbproduct.ProductQueryServiceClient
}

func NewProductQueryRepository(client pbproduct.ProductQueryServiceClient) *productQueryRepository {
	return &productQueryRepository{
		client: client,
	}
}

func (r *productQueryRepository) FindByID(ctx context.Context, product_id int) (*dto.GetProductByIDRow, error) {
	res, err := r.client.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: int32(product_id)})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &dto.GetProductByIDRow{
		ProductID:    res.Data.Id,
		MerchantID:   res.Data.MerchantId,
		CategoryID:   res.Data.CategoryId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Price:        res.Data.Price,
		CountInStock: res.Data.CountInStock,
		ImageProduct: &res.Data.ImageProduct,
	}, nil
}
