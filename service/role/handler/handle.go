package handler

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/role"
	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery       pbrole.RoleQueryServiceServer
	RoleCommand     pbrole.RoleCommandServiceServer
	UserRoleQuery   pbuserrole.UserRoleQueryServiceServer
	UserRoleCommand pbuserrole.UserRoleCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:       NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand:     NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
		UserRoleQuery:   NewUserRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		UserRoleCommand: NewUserRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
	}
}
