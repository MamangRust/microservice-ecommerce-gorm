package handler

import (
	"context"

	pbuserrole "github.com/MamangRust/microservice-ecommerce-grpc-pb/user_role"

	"github.com/MamangRust/microservice-ecommerce-grpc-role/service"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRoleCommandHandler struct {
	pbuserrole.UnimplementedUserRoleCommandServiceServer
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewUserRoleCommandHandler(roleCommand service.RoleCommandService, logger logger.LoggerInterface) pbuserrole.UserRoleCommandServiceServer {
	return &userRoleCommandHandler{
		roleCommand: roleCommand,
		logger:      logger,
	}
}

func (s *userRoleCommandHandler) AssignRoleToUser(ctx context.Context, request *pbuserrole.AssignRoleToUserRequest) (*pbuserrole.ApiResponseUserRole, error) {
	req := &requests.CreateUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbuserrole.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapToProtoUserRoleResponse(userRole),
	}, nil
}

func (s *userRoleCommandHandler) RemoveRoleFromUser(ctx context.Context, request *pbuserrole.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	req := &requests.RemoveUserRoleRequest{
		UserId: int(request.GetUserId()),
		RoleId: int(request.GetRoleId()),
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, req); err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}
