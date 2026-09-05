package repository

import (
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
)

type shippingAddressQueryRepository struct {
	client pbshipping_address.ShippingQueryServiceClient
}

func NewShippingAddressQueryRepository(client pbshipping_address.ShippingQueryServiceClient) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{
		client: client,
	}
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, order_id int) (*dto.GetShippingAddressByOrderIDRow, error) {
	res, err := r.client.FindByOrder(ctx, &pbshipping_address.FindByIdShippingRequest{Id: int32(order_id)})
	if err != nil {
		// pertahankan status gRPC dari dependency service (NotFound -> 404, dst)
		return nil, err
	}

	return &dto.GetShippingAddressByOrderIDRow{
		ShippingAddressID: res.Data.Id,
		OrderID:           res.Data.OrderId,
		Alamat:            res.Data.Alamat,
		Provinsi:          res.Data.Provinsi,
		Negara:            res.Data.Negara,
		Kota:              res.Data.Kota,
		ShippingMethod:    res.Data.ShippingMethod,
		ShippingCost:      float64(res.Data.ShippingCost),
	}, nil
}
