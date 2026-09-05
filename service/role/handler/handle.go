package handler

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	RoleQuery   pbrole.RoleQueryServiceServer
	RoleCommand pbrole.RoleCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		RoleQuery:   NewRoleQueryHandler(deps.Service.RoleQuery, deps.Logger),
		RoleCommand: NewRoleCommandHandler(deps.Service.RoleCommand, deps.Logger),
	}
}
