package repository

import (
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address"
	"context"

	dto "github.com/MamangRust/microservice-ecommerce-grpc-order/dto"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressCommandRepository struct {
	client pbshipping_address.ShippingCommandServiceClient
}

func NewShippingAddressCommandRepository(client pbshipping_address.ShippingCommandServiceClient) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{client: client}
}

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*dto.CreateShippingAddressRow, error) {
	res, err := r.client.CreateShipping(ctx, &pbshipping_address.CreateShippingAddressRequest{
		OrderId:        int32(*request.OrderID),
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int32(request.ShippingCost),
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress.WithInternal(err)
	}

	return &dto.CreateShippingAddressRow{
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

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*dto.UpdateShippingAddressRow, error) {
	var shippingID int32
	if request.ShippingID != nil {
		shippingID = int32(*request.ShippingID)
	}

	res, err := r.client.UpdateShipping(ctx, &pbshipping_address.UpdateShippingAddressRequest{
		ShippingId:     shippingID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int32(request.ShippingCost),
	})
	if err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress.WithInternal(err)
	}

	return &dto.UpdateShippingAddressRow{
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

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	_, err := r.client.DeleteShippingByOrderPermanent(ctx, &pbshipping_address.FindByIdShippingRequest{
		Id: int32(orderID),
	})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := r.client.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress.WithInternal(err)
	}

	return res.Status == "success", nil
}
