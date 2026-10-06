package repository

import (
	"context"

	roleadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/role"

	"github.com/MamangRust/microservice-ecommerce-shared/errors"
	role_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
)

type roleRepository struct {
	adapter roleadapter.QueryRepository
}

func NewRoleRepository(adapter roleadapter.QueryRepository) *roleRepository {
	return &roleRepository{
		adapter: adapter,
	}
}

func (r *roleRepository) FindByID(ctx context.Context, role_id int) (*RoleDTO, error) {
	role, err := r.adapter.FindByID(ctx, role_id)
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if role == nil {
		return nil, role_errors.ErrRoleNotFound
	}

	return &RoleDTO{
		RoleID:   role.RoleID,
		RoleName: role.RoleName,
	}, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*RoleDTO, error) {
	role, err := r.adapter.FindByName(ctx, name)
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	if role == nil {
		return nil, errors.ErrNotFound.WithMessage("role not found by name: " + name)
	}

	return &RoleDTO{
		RoleID:   role.RoleID,
		RoleName: role.RoleName,
	}, nil
}
