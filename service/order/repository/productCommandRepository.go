package repository

import (
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-order/dto"
	product_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
)

type productCommandRepository struct {
	client pbproduct.ProductCommandServiceClient
}

func NewProductCommandRepository(client pbproduct.ProductCommandServiceClient) *productCommandRepository {
	return &productCommandRepository{
		client: client,
	}
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, productID int, stock int) (*dto.UpdateProductCountStockRow, error) {
	res, err := r.client.UpdateProductCountStock(ctx, &pbproduct.UpdateProductCountStockRequest{
		ProductId: int32(productID),
		Stock:     int32(stock),
	})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &dto.UpdateProductCountStockRow{
		ProductID:    res.Data.Id,
		CountInStock: res.Data.CountInStock,
	}, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*dto.AdjustProductStockRow, error) {
	res, err := r.client.AdjustProductStock(ctx, &pbproduct.AdjustProductStockRequest{
		ProductId:   int32(productID),
		Delta:       int32(delta),
		OperationId: operationID,
	})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &dto.AdjustProductStockRow{
		ProductID:    res.Data.Id,
		CountInStock: res.Data.CountInStock,
	}, nil
}
