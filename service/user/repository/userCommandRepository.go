package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/user_errors"
	"gorm.io/gorm"
)

type userCommandRepository struct {
	db *gorm.DB
}

func NewUserCommandRepository(db *gorm.DB) *userCommandRepository {
	return &userCommandRepository{db: db}
}

func (r *userCommandRepository) Create(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	return user, nil
}

func (r *userCommandRepository) Update(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", *request.UserID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	user.Firstname = request.FirstName
	user.Lastname = request.LastName
	user.Email = request.Email
	if request.Password != "" {
		user.Password = request.Password
	}

	err = r.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Trash(ctx context.Context, user_id int) (*models.User, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("user_id = ? AND deleted_at IS NULL", user_id).
		Update("deleted_at", now)
	if result.Error != nil {
		return nil, user_errors.ErrTrashedUser.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, user_errors.ErrUserNotFound
	}

	var user models.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) Restore(ctx context.Context, user_id int) (*models.User, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).
		Where("user_id = ? AND deleted_at IS NOT NULL", user_id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return nil, user_errors.ErrRestoreUser.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, user_errors.ErrUserNotFound
	}

	var user models.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&user).Error; err != nil {
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) DeletePermanent(ctx context.Context, user_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL", user_id).
		Delete(&models.User{})
	if result.Error != nil {
		return false, user_errors.ErrDeleteUserPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, user_errors.ErrUserNotFound
	}
	return true, nil
}

func (r *userCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).
		Where("deleted_at IS NOT NULL").Update("deleted_at", nil).Error
	if err != nil {
		return false, user_errors.ErrRestoreAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").Delete(&models.User{}).Error
	if err != nil {
		return false, user_errors.ErrDeleteAllUsers.WithInternal(err)
	}
	return true, nil
}

func (r *userCommandRepository) UpdateIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", user_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	user.IsVerified = &is_verified
	err = r.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

func (r *userCommandRepository) UpdatePassword(ctx context.Context, user_id int, password string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", user_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	user.Password = password
	err = r.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}
	return &user, nil
}

// Ensure shared_errors is used
var _ = shared_errors.ErrInternal
