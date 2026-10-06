// Package role is the shared gRPC adapter for the role domain. It owns the only
// place that talks to pb_role's RoleQueryService.
package role

import (
	"context"
	"time"

	pb_role "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
)

// Role is the transport-neutral view of a role row.
type Role struct {
	RoleID    int32
	RoleName  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// QueryRepository is the read path consumers use to resolve roles.
type QueryRepository interface {
	FindByID(ctx context.Context, id int) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	FindAll(ctx context.Context, search string, page, pageSize int) ([]Role, int, error)
}

// Repository implements QueryRepository.
type Repository struct {
	query pb_role.RoleQueryServiceClient
	guard *resilience.DependencyGuard
}

// New builds the role adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(query pb_role.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter is the pre-New constructor; it forwards to New.
func NewQueryAdapter(query pb_role.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, opts...)
}

// NewAdapter is the pre-New constructor; it forwards to New.
func NewAdapter(query pb_role.RoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return New(query, opts...)
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

// FindByID implements QueryRepository.
func (r *Repository) FindByID(ctx context.Context, id int) (*Role, error) {
	var res *pb_role.ApiResponseRole
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindByIdRole(callCtx, &pb_role.FindByIdRoleRequest{RoleId: int32(id)})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toRole(res.Data), nil
}

// FindByName implements QueryRepository.
func (r *Repository) FindByName(ctx context.Context, name string) (*Role, error) {
	var res *pb_role.ApiResponseRole
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindByNameRole(callCtx, &pb_role.FindByNameRoleRequest{Name: name})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toRole(res.Data), nil
}

// FindAll implements QueryRepository.
func (r *Repository) FindAll(ctx context.Context, search string, page, pageSize int) ([]Role, int, error) {
	var res *pb_role.ApiResponsePaginationRole
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindAllRole(callCtx, &pb_role.FindAllRoleRequest{
			Search:   search,
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
		return callErr
	})
	if err != nil {
		return nil, 0, err
	}

	roles := make([]Role, 0, len(res.Data))
	for _, item := range res.Data {
		roles = append(roles, *toRole(item))
	}

	total := 0
	if res.Pagination != nil {
		total = int(res.Pagination.TotalRecords)
	}
	return roles, total, nil
}

func toRole(r *pb_role.RoleResponse) *Role {
	if r == nil {
		return nil
	}
	return &Role{
		RoleID:    r.Id,
		RoleName:  r.Name,
		CreatedAt: adapter.ParseTimePtr(r.CreatedAt),
		UpdatedAt: adapter.ParseTimePtr(r.UpdatedAt),
	}
}

var (
	_ QueryRepository = (*Repository)(nil)
)
