// Package user is the shared gRPC adapter for the user domain. It owns the only
// place that talks to pb_user's UserQueryService/UserCommandService.
package user

import (
	"context"
	"time"

	pb_user "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	"github.com/MamangRust/microservice-ecommerce-pkg/resilience"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// User is the transport-neutral view of a user row. Password is only populated
// by the calls whose RPC returns it (FindByEmail, Create, Update*).
type User struct {
	UserID    int32
	Firstname string
	Lastname  string
	Email     string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// QueryRepository is the read path consumers use to resolve users.
type QueryRepository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*User, error)
	FindByVerificationCode(ctx context.Context, code string) (*User, error)
}

// CommandRepository is the write path consumers use to create/update users.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.RegisterRequest) (*User, error)
	UpdateIsVerified(ctx context.Context, id int, isVerified bool) (*User, error)
	UpdatePassword(ctx context.Context, id int, password string) (*User, error)
}

// Repository implements QueryRepository and (when a command client is supplied)
// CommandRepository.
type Repository struct {
	query   pb_user.UserQueryServiceClient
	command pb_user.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds the user adapter from its query and command clients. Passing zero
// options leaves the guard nil, which makes DependencyGuard.Call a plain
// passthrough.
func New(query pb_user.UserQueryServiceClient, command pb_user.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

// NewQueryAdapter builds a read-only user adapter.
func NewQueryAdapter(query pb_user.UserQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, nil, opts)
}

// NewAdapter is the pre-New constructor; it forwards to New.
func NewAdapter(query pb_user.UserQueryServiceClient, command pb_user.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	return newRepository(query, command, opts)
}

func newRepository(query pb_user.UserQueryServiceClient, command pb_user.UserCommandServiceClient, opts []adapter.GuardOption) *Repository {
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
func (r *Repository) FindByID(ctx context.Context, id int) (*User, error) {
	var res *pb_user.ApiResponseUser
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindById(callCtx, &pb_user.FindByIdUserRequest{Id: int32(id)})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toUser(res.Data), nil
}

// FindByEmail implements QueryRepository. The returned user carries the
// password hash so the auth service can verify credentials.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var res *pb_user.ApiResponseUserWithPassword
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindByEmail(callCtx, &pb_user.FindByEmailRequest{Email: email})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return &User{
		UserID:    res.Data.Id,
		Firstname: res.Data.Firstname,
		Lastname:  res.Data.Lastname,
		Email:     res.Data.Email,
		Password:  res.Data.Password,
	}, nil
}

// FindByEmailAndVerify implements QueryRepository. It resolves the account by
// email; verification status is not part of the lookup.
func (r *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*User, error) {
	return r.FindByEmail(ctx, email)
}

// FindByVerificationCode implements QueryRepository.
func (r *Repository) FindByVerificationCode(ctx context.Context, code string) (*User, error) {
	var res *pb_user.ApiResponseUser
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.query.FindByVerificationCode(callCtx, &pb_user.FindByVerificationCodeRequest{VerificationCode: code})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toUser(res.Data), nil
}

// Create implements CommandRepository.
func (r *Repository) Create(ctx context.Context, req *requests.RegisterRequest) (*User, error) {
	var res *pb_user.ApiResponseUser
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.Create(callCtx, &pb_user.CreateUserRequest{
			Firstname:       req.FirstName,
			Lastname:        req.LastName,
			Email:           req.Email,
			Password:        req.Password,
			ConfirmPassword: req.ConfirmPassword,
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toUser(res.Data), nil
}

// UpdateIsVerified implements CommandRepository.
func (r *Repository) UpdateIsVerified(ctx context.Context, id int, isVerified bool) (*User, error) {
	var res *pb_user.ApiResponseUser
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.UpdateIsVerified(callCtx, &pb_user.UpdateUserIsVerifiedRequest{
			Id:         int32(id),
			IsVerified: isVerified,
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toUser(res.Data), nil
}

// UpdatePassword implements CommandRepository.
func (r *Repository) UpdatePassword(ctx context.Context, id int, password string) (*User, error) {
	var res *pb_user.ApiResponseUser
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.command.UpdatePassword(callCtx, &pb_user.UpdateUserPasswordRequest{
			Id:       int32(id),
			Password: password,
		})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	return toUser(res.Data), nil
}

func toUser(u *pb_user.UserResponse) *User {
	if u == nil {
		return nil
	}
	return &User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		CreatedAt: adapter.ParseTimePtr(u.CreatedAt),
		UpdatedAt: adapter.ParseTimePtr(u.UpdatedAt),
	}
}

var (
	_ QueryRepository   = (*Repository)(nil)
	_ CommandRepository = (*Repository)(nil)
)
