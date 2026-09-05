package repository

import (
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-order/dto"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	order_item_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderItemCommandRepository struct {
	client pborder_item.OrderItemCommandServiceClient
}

func NewOrderItemCommandRepository(client pborder_item.OrderItemCommandServiceClient) *orderItemCommandRepository {
	return &orderItemCommandRepository{client: client}
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*dto.CreateOrderItemRow, error) {
	res, err := r.client.CreateOrderItem(ctx, &pborder_item.CreateOrderItemRecordRequest{
		OrderId:   int32(req.OrderID),
		ProductId: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})
	if err != nil {
		return nil, order_item_errors.ErrCreateOrderItem.WithInternal(err)
	}

	return &dto.CreateOrderItemRow{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*dto.UpdateOrderItemRow, error) {
	res, err := r.client.UpdateOrderItem(ctx, &pborder_item.UpdateOrderItemRecordRequest{
		OrderItemId: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})
	if err != nil {
		return nil, order_item_errors.ErrUpdateOrderItem.WithInternal(err)
	}

	return &dto.UpdateOrderItemRow{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, orderID int) (*dto.OrderItem, error) {
	res, err := r.client.TrashOrderItem(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	if err != nil {
		return nil, order_item_errors.ErrTrashedOrderItem.WithInternal(err)
	}

	return &dto.OrderItem{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, orderID int) (*dto.OrderItem, error) {
	res, err := r.client.RestoreOrderItem(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	if err != nil {
		return nil, order_item_errors.ErrRestoreOrderItem.WithInternal(err)
	}

	return &dto.OrderItem{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := r.client.DeleteOrderItemPermanent(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	res, err := r.client.DeleteOrderItemByOrderPermanent(ctx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	res, err := r.client.RestoreAllOrdersItem(ctx, &emptypb.Empty{})
	if err != nil {
		return false, order_item_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := r.client.DeleteAllPermanentOrdersItem(ctx, &emptypb.Empty{})
	if err != nil {
		return false, order_item_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}
