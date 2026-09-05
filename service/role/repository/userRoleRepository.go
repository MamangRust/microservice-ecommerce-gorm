package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	userRole := &models.UserRole{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	}
	err := r.db.WithContext(ctx).Create(userRole).Error
	if err != nil {
		return nil, role_errors.ErrAssignRole.WithInternal(err)
	}
	return userRole, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", req.UserId, req.RoleId).
		Delete(&models.UserRole{})
	if result.Error != nil {
		return role_errors.ErrRemoveRole.WithInternal(result.Error)
	}
	return nil
}
