package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) *userQueryRepository {
	return &userQueryRepository{db: db}
}

func (r *userQueryRepository) FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult

	query := `
		SELECT user_id, firstname, lastname, email, password, created_at, updated_at,
			COUNT(*) OVER() AS total_count
		FROM users
		WHERE deleted_at IS NULL
			AND (? = '' OR firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindByID(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", user_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByIDWithPassword(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", user_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult

	query := `
		SELECT user_id, firstname, lastname, email, password, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM users
		WHERE deleted_at IS NULL
			AND (? = '' OR firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, user_errors.ErrFindActiveUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*UserResult

	query := `
		SELECT user_id, firstname, lastname, email, password, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM users
		WHERE deleted_at IS NOT NULL
			AND (? = '' OR firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}
	return results, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("verification_code = ?", code).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &user, nil
}
