package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/role_errors"
	"gorm.io/gorm"
)

type roleQueryRepository struct {
	db *gorm.DB
}

func NewRoleQueryRepository(db *gorm.DB) *roleQueryRepository {
	return &roleQueryRepository{db: db}
}

func (r *roleQueryRepository) FindAll(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*RoleResult

	query := `
		SELECT role_id, role_name, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM roles
		WHERE deleted_at IS NULL
			AND (? = '' OR role_name ILIKE ?)
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, role_errors.ErrFindAllRoles.WithInternal(err)
	}

	return results, nil
}

func (r *roleQueryRepository) FindByID(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("role_id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("role_name = ? AND deleted_at IS NULL", name).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, userID int) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.role_id").
		Where("user_roles.user_id = ? AND roles.deleted_at IS NULL", userID).
		Find(&roles).Error
	if err != nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roles, nil
}

func (r *roleQueryRepository) FindActive(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*RoleResult

	query := `
		SELECT role_id, role_name, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM roles
		WHERE deleted_at IS NULL
			AND (? = '' OR role_name ILIKE ?)
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, role_errors.ErrFindActiveRoles.WithInternal(err)
	}

	return results, nil
}

func (r *roleQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllRole) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*RoleResult

	query := `
		SELECT role_id, role_name, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM roles
		WHERE deleted_at IS NOT NULL
			AND (? = '' OR role_name ILIKE ?)
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, role_errors.ErrFindTrashedRoles.WithInternal(err)
	}

	return results, nil
}

// Ensure unused imports are referenced
var _ = time.Now
