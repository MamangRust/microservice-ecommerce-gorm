package repository

import (
	pbroles "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"gorm.io/gorm"
)

// GuardOptions carries the resilience guard options for the role and user-role
// dependencies. UserRole is hosted by the role service.
type GuardOptions struct {
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
	UserRole    UserRoleRepository
}

type Deps struct {
	Db                    *gorm.DB
	RoleQueryClient       pbroles.RoleQueryServiceClient
	UserRoleQueryClient   pbuserrole.UserRoleQueryServiceClient
	UserRoleCommandClient pbuserrole.UserRoleCommandServiceClient
	Guard                 GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.Db),
		UserQuery:   NewUserQueryRepository(deps.Db),
		Role:        NewRoleRepository(roleadapter.New(deps.RoleQueryClient, deps.Guard.Role...)),
		UserRole: userroleadapter.New(
			deps.UserRoleQueryClient,
			deps.UserRoleCommandClient,
			deps.Guard.UserRole...,
		),
	}
}
