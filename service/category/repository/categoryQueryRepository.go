package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
	"gorm.io/gorm"
)

type categoryQueryRepository struct {
	db *gorm.DB
}

func NewCategoryQueryRepository(db *gorm.DB) *categoryQueryRepository {
	return &categoryQueryRepository{db: db}
}

func (r *categoryQueryRepository) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*CategoryResult

	query := `
		SELECT category_id, name, description, slug_category, image_category, created_at, updated_at,
			COUNT(*) OVER() AS total_count
		FROM categories
		WHERE deleted_at IS NULL
			AND (? = '' OR name ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*CategoryResult

	query := `
		SELECT category_id, name, description, slug_category, image_category, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM categories
		WHERE deleted_at IS NULL
			AND (? = '' OR name ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*CategoryResult

	query := `
		SELECT category_id, name, description, slug_category, image_category, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM categories
		WHERE deleted_at IS NOT NULL
			AND (? = '' OR name ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}
	return results, nil
}

func (r *categoryQueryRepository) FindByID(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Where("category_id = ? AND deleted_at IS NULL", category_id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
		}
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}
	return &category, nil
}

func (r *categoryQueryRepository) FindByIDTrashed(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).Unscoped().Where("category_id = ?", category_id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
		}
		return nil, category_errors.ErrFindCategoryByIdTrashed.WithInternal(err)
	}
	return &category, nil
}
