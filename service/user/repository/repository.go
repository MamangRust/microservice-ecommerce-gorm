package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	"gorm.io/gorm"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
}

func NewRepositories(DB *gorm.DB, roleClient pbrole.RoleQueryServiceClient) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(DB),
		UserQuery:   NewUserQueryRepository(DB),
		Role:        NewRoleRepository(roleClient),
	}
}
