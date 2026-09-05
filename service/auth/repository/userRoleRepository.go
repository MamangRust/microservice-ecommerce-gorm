package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	userrole_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/user_role_errors"
)

type userRoleRepository struct {
	client pbrole.RoleCommandServiceClient
}

func NewUserRoleRepository(client pbrole.RoleCommandServiceClient) UserRoleRepository {
	return &userRoleRepository{
		client: client,
	}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*AuthUserRole, error) {
	protoReq := &pbrole.AssignRoleToUserRequest{
		UserId: int32(req.UserId),
		RoleId: int32(req.RoleId),
	}

	res, err := r.client.AssignRoleToUser(ctx, protoReq)
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}

	return &AuthUserRole{
		UserRoleID: res.Data.UserRoleId,
		UserID:     res.Data.UserId,
		RoleID:     res.Data.RoleId,
	}, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	protoReq := &pbrole.RemoveRoleFromUserRequest{
		UserId: int32(req.UserId),
		RoleId: int32(req.RoleId),
	}

	_, err := r.client.RemoveRoleFromUser(ctx, protoReq)
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}

	return nil
}
