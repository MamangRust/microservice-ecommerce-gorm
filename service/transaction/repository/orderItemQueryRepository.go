package repository

import (
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
)

type orderItemRepository struct {
	client pborder_item.OrderItemQueryServiceClient
}

func NewOrderItemRepository(client pborder_item.OrderItemQueryServiceClient) *orderItemRepository {
	return &orderItemRepository{
		client: client,
	}
}

func (r *orderItemRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*dto.GetOrderItemsByOrderRow, error) {
	res, err := r.client.FindOrderItemByOrder(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
	}

	var items []*dto.GetOrderItemsByOrderRow
	for _, item := range res.Data {
		items = append(items, &dto.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    int32(item.Quantity),
			Price:       int32(item.Price),
		})
	}

	return items, nil
}
