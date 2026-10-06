// Package user_role is the shared gRPC adapter for the user-role domain. It owns
// the only place that talks to pb_user_role's UserRoleQueryService and
// UserRoleCommandService.
package user_role

import (
	"context"

	pb_user_role "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// UserRole is the transport-neutral view of a user-role assignment.
type UserRole struct {
	UserRoleID int32
	UserID     int32
	RoleID     int32
}

// Role is the transport-neutral view of a role returned by FindByUserId.
type Role struct {
	RoleID   int32
	RoleName string
}

// QueryRepository is the read path consumers use to resolve a user's roles.
type QueryRepository interface {
	FindByUserId(ctx context.Context, userID int) ([]Role, error)
}

// CommandRepository is the write path consumers use to (un)assign roles.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements QueryRepository and (when a command client is supplied)
// CommandRepository.
type Repository struct {
	query   pb_user_role.UserRoleQueryServiceClient
	command pb_user_role.UserRoleCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the user-role adapter from its query and command clients. Passing
// zero options leaves the guard nil, which makes DependencyGuard.Call a plain
// passthrough.
func New(query pb_user_role.UserRoleQueryServiceClient, command pb_user_role.UserRoleCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

// NewQueryAdapter builds a read-only user-role adapter.
func NewQueryAdapter(query pb_user_role.UserRoleQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, nil, opts)
}

// NewCommandAdapter builds a write-only user-role adapter.
func NewCommandAdapter(command pb_user_role.UserRoleCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(nil, command, opts)
}

// NewAdapter is the pre-New constructor; it forwards to New.
func NewAdapter(query pb_user_role.UserRoleQueryServiceClient, command pb_user_role.UserRoleCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

func newRepository(query pb_user_role.UserRoleQueryServiceClient, command pb_user_role.UserRoleCommandServiceClient, opts []adapter.GuardOption) *Repository {
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

// FindByUserId implements QueryRepository.
func (r *Repository) FindByUserId(ctx context.Context, userID int) ([]Role, error) {
	return adapter.Call(r.guard, ctx, func(callCtx context.Context) ([]Role, error) {
		res, err := r.query.FindByUserId(callCtx, &pb_user_role.FindByIdUserRoleRequest{UserId: int32(userID)})
		if err != nil {
			return nil, err
		}
		roles := make([]Role, 0, len(res.Data))
		for _, item := range res.Data {
			if item == nil {
				continue
			}
			roles = append(roles, Role{RoleID: item.Id, RoleName: item.Name})
		}
		return roles, nil
	})
}

// AssignRoleToUser implements CommandRepository.
func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*UserRole, error) {
	return adapter.Call(r.guard, ctx, func(callCtx context.Context) (*UserRole, error) {
		res, err := r.command.AssignRoleToUser(callCtx, &pb_user_role.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		if err != nil {
			return nil, err
		}
		if res.Data == nil {
			return &UserRole{UserID: int32(req.UserId), RoleID: int32(req.RoleId)}, nil
		}
		return &UserRole{
			UserRoleID: res.Data.UserRoleId,
			UserID:     res.Data.UserId,
			RoleID:     res.Data.RoleId,
		}, nil
	})
}

// RemoveRoleFromUser implements CommandRepository.
func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	return r.guard.Call(ctx, func(callCtx context.Context) error {
		_, err := r.command.RemoveRoleFromUser(callCtx, &pb_user_role.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return err
	})
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
