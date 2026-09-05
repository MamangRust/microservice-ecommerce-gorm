package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
	"gorm.io/gorm"
)

type categoryCommandRepository struct {
	db *gorm.DB
}

func NewCategoryCommandRepository(db *gorm.DB) *categoryCommandRepository {
	return &categoryCommandRepository{db: db}
}

func (r *categoryCommandRepository) Create(ctx context.Context, request *requests.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:          request.Name,
		Description:   &request.Description,
		SlugCategory:  request.SlugCategory,
		ImageCategory: &request.ImageCategory,
	}
	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}
	return category, nil
}

func (r *categoryCommandRepository) Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Where("category_id = ?", *request.CategoryID).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category_errors.ErrCategoryNotFound
		}
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}

	category.Name = request.Name
	category.Description = &request.Description
	category.SlugCategory = request.SlugCategory
	category.ImageCategory = &request.ImageCategory

	err = r.db.WithContext(ctx).Save(&category).Error
	if err != nil {
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Trash(ctx context.Context, category_id int) (*models.Category, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.Category{}).
		Where("category_id = ? AND deleted_at IS NULL", category_id).
		Update("deleted_at", now)
	if result.Error != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}

	var category models.Category
	if err := r.db.WithContext(ctx).Where("category_id = ?", category_id).First(&category).Error; err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) Restore(ctx context.Context, category_id int) (*models.Category, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Category{}).
		Where("category_id = ? AND deleted_at IS NOT NULL", category_id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, category_errors.ErrCategoryNotFound
	}

	var category models.Category
	if err := r.db.WithContext(ctx).Where("category_id = ?", category_id).First(&category).Error; err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryCommandRepository) DeletePermanent(ctx context.Context, category_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("category_id = ? AND deleted_at IS NOT NULL", category_id).
		Delete(&models.Category{})
	if result.Error != nil {
		return false, category_errors.ErrDeleteCategoryPermanently.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, category_errors.ErrCategoryNotFound
	}
	return true, nil
}

func (r *categoryCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().Model(&models.Category{}).
		Where("deleted_at IS NOT NULL").Update("deleted_at", nil).Error
	if err != nil {
		return false, category_errors.ErrRestoreAllCategories.WithInternal(err)
	}
	return true, nil
}

func (r *categoryCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").Delete(&models.Category{}).Error
	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories.WithInternal(err)
	}
	return true, nil
}

var _ = shared_errors.ErrInternal
