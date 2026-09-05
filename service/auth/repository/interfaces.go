package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"gorm.io/gorm"
)

// AuthUser mirrors the user fields the auth service needs from the user service gRPC.
type AuthUser struct {
	UserID    int32
	Email     string
	Firstname string
	Lastname  string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*AuthUser, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*AuthUser, error)
	FindById(ctx context.Context, user_id int) (*AuthUser, error)
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*AuthUser, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*AuthUser, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*AuthUser, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*AuthUser, error)
}

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, code string) (*models.ResetToken, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	CreateResetTokenInTx(ctx context.Context, tx *gorm.DB, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

// AuthUserRole mirrors the user role fields.
type AuthUserRole struct {
	UserRoleID int32
	UserID     int32
	RoleID     int32
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*AuthUserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// AuthRole mirrors the role fields.
type AuthRole struct {
	RoleID   int32
	RoleName string
}

type RoleRepository interface {
	FindById(ctx context.Context, id int) (*AuthRole, error)
	FindByName(ctx context.Context, name string) (*AuthRole, error)
}

// Repositories aggregates all auth repositories.
type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}
