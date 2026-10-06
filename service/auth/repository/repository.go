package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user_role"
	"gorm.io/gorm"
)

type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
}

type GuardOptions struct {
	User     []adapter.GuardOption
	UserRole []adapter.GuardOption
	Role     []adapter.GuardOption
}

var (
	_ UserRepository     = (*useradapter.Repository)(nil)
	_ RoleRepository     = (*roleadapter.Repository)(nil)
	_ UserRoleRepository = (*userroleadapter.Repository)(nil)
)

type Deps struct {
	Db              *gorm.DB
	User            pbuser.UserQueryServiceClient
	UserCommand     pbuser.UserCommandServiceClient
	Role            pbrole.RoleQueryServiceClient
	UserRoleCommand pbuserrole.UserRoleCommandServiceClient
	Guards          GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		User:         useradapter.New(deps.User, deps.UserCommand, deps.Guards.User...),
		RefreshToken: NewRefreshTokenRepository(deps.Db),
		UserRole:     userroleadapter.NewCommandAdapter(deps.UserRoleCommand, deps.Guards.UserRole...),
		Role:         roleadapter.New(deps.Role, deps.Guards.Role...),
		ResetToken:   NewResetTokenRepository(deps.Db),
	}
}
