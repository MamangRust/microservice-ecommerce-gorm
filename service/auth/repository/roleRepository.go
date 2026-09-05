package repository

import (
	pbrole "github.com/MamangRust/microservice-ecommerce-grpc/pb/role"
	"context"
	"fmt"

)

type roleRepository struct {
	client pbrole.RoleQueryServiceClient
}

func NewRoleRepository(client pbrole.RoleQueryServiceClient) RoleRepository {
	return &roleRepository{
		client: client,
	}
}

func (r *roleRepository) FindById(ctx context.Context, id int) (*AuthRole, error) {
	res, err := r.client.FindByIdRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by ID %d: %w", id, err)
	}

	return &AuthRole{
		RoleID:   res.Data.Id,
		RoleName: res.Data.Name,
	}, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*AuthRole, error) {
	res, err := r.client.FindByNameRole(ctx, &pbrole.FindByNameRoleRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find role by name %s: %w", name, err)
	}

	return &AuthRole{
		RoleID:   res.Data.Id,
		RoleName: res.Data.Name,
	}, nil
}
