// Package shippingaddress is the shared gRPC adapter for the shipping-address
// domain. It owns the only place that talks to pb_shipping_address's
// ShippingQueryService/ShippingCommandService.
package shippingaddress

import (
	"context"

	pb_shipping "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ShippingAddress is the transport-neutral view of a shipping-address row.
type ShippingAddress struct {
	ShippingAddressID int32
	OrderID           int32
	Alamat            string
	Provinsi          string
	Kota              string
	Negara            string
	Courier           string
	ShippingMethod    string
	ShippingCost      int32
}

// QueryRepository is the read path consumers use to resolve a shipping address.
type QueryRepository interface {
	FindByID(ctx context.Context, id int) (*ShippingAddress, error)
	FindByOrder(ctx context.Context, orderID int) (*ShippingAddress, error)
}

// CommandRepository is the write path consumers use to manage shipping
// addresses.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*ShippingAddress, error)
	Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*ShippingAddress, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements QueryRepository and (when a command client is supplied)
// CommandRepository.
type Repository struct {
	query   pb_shipping.ShippingQueryServiceClient
	command pb_shipping.ShippingCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the shipping-address adapter from its query and command clients.
// Passing zero options leaves the guard nil, which makes DependencyGuard.Call a
// plain passthrough.
func New(query pb_shipping.ShippingQueryServiceClient, command pb_shipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

// NewQueryAdapter builds a read-only shipping-address adapter.
func NewQueryAdapter(query pb_shipping.ShippingQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, nil, opts)
}

// NewAdapter is the pre-New constructor; it forwards to New.
func NewAdapter(query pb_shipping.ShippingQueryServiceClient, command pb_shipping.ShippingCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

func newRepository(query pb_shipping.ShippingQueryServiceClient, command pb_shipping.ShippingCommandServiceClient, opts []adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindByID implements QueryRepository.
func (r *Repository) FindByID(ctx context.Context, id int) (*ShippingAddress, error) {
	var res *pb_shipping.ApiResponseShipping
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindById(callCtx, &pb_shipping.FindByIdShippingRequest{Id: int32(id)})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toShippingAddress(res.Data), nil
}

// FindByOrder implements QueryRepository.
func (r *Repository) FindByOrder(ctx context.Context, orderID int) (*ShippingAddress, error) {
	var res *pb_shipping.ApiResponseShipping
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindByOrder(callCtx, &pb_shipping.FindByIdShippingRequest{Id: int32(orderID)})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toShippingAddress(res.Data), nil
}

// Create implements CommandRepository.
func (r *Repository) Create(ctx context.Context, req *requests.CreateShippingAddressRequest) (*ShippingAddress, error) {
	var res *pb_shipping.ApiResponseShipping
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.CreateShipping(callCtx, &pb_shipping.CreateShippingAddressRequest{
			OrderId:        int32(*req.OrderID),
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toShippingAddress(res.Data), nil
}

// Update implements CommandRepository.
func (r *Repository) Update(ctx context.Context, req *requests.UpdateShippingAddressRequest) (*ShippingAddress, error) {
	var shippingID int32
	if req.ShippingID != nil {
		shippingID = int32(*req.ShippingID)
	}

	var res *pb_shipping.ApiResponseShipping
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.UpdateShipping(callCtx, &pb_shipping.UpdateShippingAddressRequest{
			ShippingId:     shippingID,
			Alamat:         req.Alamat,
			Provinsi:       req.Provinsi,
			Kota:           req.Kota,
			Negara:         req.Negara,
			Courier:        req.Courier,
			ShippingMethod: req.ShippingMethod,
			ShippingCost:   int32(req.ShippingCost),
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toShippingAddress(res.Data), nil
}

// DeleteByOrderIDPermanent implements CommandRepository.
func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		_, callErr := r.command.DeleteShippingByOrderPermanent(callCtx, &pb_shipping.FindByIdShippingRequest{
			Id: int32(orderID),
		})
		return callErr
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteAll implements CommandRepository.
func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	var res *pb_shipping.ApiResponseShippingAll
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.DeleteAllShippingPermanent(callCtx, &emptypb.Empty{})
		return callErr
	})
	if err != nil {
		return false, err
	}
	return res.Status == "success", nil
}

func toShippingAddress(s *pb_shipping.ShippingResponse) *ShippingAddress {
	if s == nil {
		return nil
	}
	return &ShippingAddress{
		ShippingAddressID: s.Id,
		OrderID:           s.OrderId,
		Alamat:            s.Alamat,
		Provinsi:          s.Provinsi,
		Kota:              s.Kota,
		Negara:            s.Negara,
		Courier:           s.Courier,
		ShippingMethod:    s.ShippingMethod,
		ShippingCost:      s.ShippingCost,
	}
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
