package repository

import (
	"context"

	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"gorm.io/gorm"
)

// AuthUser mirrors the user fields the auth service needs from the user service
// gRPC. It is the shared user adapter's view.
type AuthUser = useradapter.User

// AuthRole mirrors the role fields.
type AuthRole = roleadapter.Role

// AuthUserRole mirrors the user role fields.
type AuthUserRole = userroleadapter.UserRole

// UserRepository is the user query+command surface, provided by the shared user
// gRPC adapter so this service never holds a raw gRPC client.
type UserRepository interface {
	useradapter.QueryRepository
	useradapter.CommandRepository
}

// RoleRepository is the shared role query contract.
type RoleRepository = roleadapter.QueryRepository

// UserRoleRepository is the shared user-role command contract.
type UserRoleRepository = userroleadapter.CommandRepository

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
