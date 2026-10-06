// Package merchant is the shared gRPC adapter for the merchant domain. It owns
// the only place that talks to pb_merchant's MerchantQueryService, so the seven
// services that need merchant data no longer hold a raw *ServiceClient.
package merchant

import (
	"context"
	"time"

	pb_merchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
)

// Merchant is the transport-neutral view of a merchant row.
type Merchant struct {
	MerchantID   int32
	UserID       int32
	Name         string
	Description  string
	Address      string
	ContactEmail string
	ContactPhone string
	Status       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

// QueryRepository is the read path consumers use to resolve a merchant by id.
type QueryRepository interface {
	FindByID(ctx context.Context, id int) (*Merchant, error)
}

// Repository implements QueryRepository on top of the merchant query service.
type Repository struct {
	query pb_merchant.MerchantQueryServiceClient
	guard *resilience.DependencyGuard
}

// New builds the merchant query adapter. Passing zero options leaves the guard
// nil, which makes DependencyGuard.Call a plain passthrough.
func New(query pb_merchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter is the pre-New constructor; it forwards to New.
func NewQueryAdapter(query pb_merchant.MerchantQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, opts...)
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindByID implements QueryRepository. The returned error is the raw gRPC
// error; each consumer keeps its own error mapping.
func (r *Repository) FindByID(ctx context.Context, id int) (*Merchant, error) {
	var res *pb_merchant.ApiResponseMerchant
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindById(callCtx, &pb_merchant.FindByIdMerchantRequest{Id: int32(id)})
		return callErr
	})
	if err != nil {
		return nil, err
	}

	return &Merchant{
		MerchantID:   res.Data.Id,
		UserID:       res.Data.UserId,
		Name:         res.Data.Name,
		Description:  res.Data.Description,
		Address:      res.Data.Address,
		ContactEmail: res.Data.ContactEmail,
		ContactPhone: res.Data.ContactPhone,
		Status:       res.Data.Status,
		CreatedAt:    adapter.ParseTimePtr(res.Data.CreatedAt),
		UpdatedAt:    adapter.ParseTimePtr(res.Data.UpdatedAt),
	}, nil
}

// compile-time check that Repository satisfies the interface consumers use.
var _ QueryRepository = (*Repository)(nil)
