package repository

import (
	pborder "github.com/MamangRust/microservice-ecommerce-grpc/pb/order"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
)

type orderQueryRepository struct {
	client pborder.OrderQueryServiceClient
}

func NewOrderQueryRepository(client pborder.OrderQueryServiceClient) *orderQueryRepository {
	return &orderQueryRepository{
		client: client,
	}
}

func (r *orderQueryRepository) FindByID(ctx context.Context, order_id int) (*dto.GetOrderByIDRow, error) {
	res, err := r.client.FindById(ctx, &pborder.FindByIdOrderRequest{Id: int32(order_id)})
	if err != nil {
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
	}

	return &dto.GetOrderByIDRow{
		OrderID:    res.Data.Id,
		UserID:     res.Data.UserId,
		MerchantID: res.Data.MerchantId,
		TotalPrice: res.Data.TotalPrice,
	}, nil
}
