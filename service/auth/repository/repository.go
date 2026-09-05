package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

func NewRepositories(DB *gorm.DB,
	userQuery pbuser.UserQueryServiceClient,
	userCommand pbuser.UserCommandServiceClient,
	roleQuery pbrole.RoleQueryServiceClient,
	roleCommand pbrole.RoleCommandServiceClient,
) *Repositories {
	return &Repositories{
		User:         NewUserRepository(userQuery, userCommand),
		RefreshToken: NewRefreshTokenRepository(DB),
		UserRole:     NewUserRoleRepository(roleCommand),
		Role:         NewRoleRepository(roleQuery),
		ResetToken:   NewResetTokenRepository(DB),
	}
}
