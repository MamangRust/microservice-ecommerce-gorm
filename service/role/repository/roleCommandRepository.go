package repository

import (
	"context"
	"strings"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"gorm.io/gorm"
)

type roleCommandRepository struct {
	db *gorm.DB
}

func NewRoleCommandRepository(db *gorm.DB) *roleCommandRepository {
	return &roleCommandRepository{db: db}
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "UNIQUE constraint")
}

func (r *roleCommandRepository) Create(ctx context.Context, req *requests.CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{RoleName: req.Name}
	err := r.db.WithContext(ctx).Create(role).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}
	return role, nil
}

func (r *roleCommandRepository) Update(ctx context.Context, req *requests.UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("role_id = ?", *req.ID).First(&role).Error
	if err != nil {
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}

	role.RoleName = req.Name
	err = r.db.WithContext(ctx).Save(&role).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil, role_errors.ErrRoleConflict
		}
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}

	return &role, nil
}

func (r *roleCommandRepository) Trash(ctx context.Context, id int) (*models.Role, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.Role{}).
		Where("role_id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now)
	if result.Error != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, role_errors.ErrRoleNotFound
	}

	var role models.Role
	if err := r.db.WithContext(ctx).Where("role_id = ?", id).First(&role).Error; err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) Restore(ctx context.Context, id int) (*models.Role, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Role{}).
		Where("role_id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, role_errors.ErrRoleNotFound
	}

	var role models.Role
	if err := r.db.WithContext(ctx).Where("role_id = ?", id).First(&role).Error; err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}
	return &role, nil
}

func (r *roleCommandRepository) DeletePermanent(ctx context.Context, roleID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("role_id = ? AND deleted_at IS NOT NULL", roleID).
		Delete(&models.Role{})
	if result.Error != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, role_errors.ErrRoleNotFound
	}
	return true, nil
}

func (r *roleCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().Model(&models.Role{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil).Error
	if err != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(err)
	}
	return true, nil
}

func (r *roleCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Role{}).Error
	if err != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(err)
	}
	return true, nil
}
