package repository

import (
	"context"

	categoryadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/category"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
	"gorm.io/gorm"
)

type productQueryRepository struct {
	db           *gorm.DB
	categoryRepo categoryadapter.QueryRepository
}

func NewProductQueryRepository(db *gorm.DB, categoryRepo categoryadapter.QueryRepository) *productQueryRepository {
	return &productQueryRepository{db: db, categoryRepo: categoryRepo}
}

func (r *productQueryRepository) FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at,
			COUNT(*) OVER() as total_count
		FROM products p
		WHERE p.deleted_at IS NULL
			AND (? = '' OR p.name ILIKE ? OR p.description ILIKE ?)
		ORDER BY p.product_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, product_errors.ErrFindAllProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at,
			COUNT(*) OVER() as total_count
		FROM products p
		WHERE p.deleted_at IS NULL
			AND (? = '' OR p.name ILIKE ? OR p.description ILIKE ?)
		ORDER BY p.product_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, product_errors.ErrFindActiveProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at,
			COUNT(*) OVER() as total_count
		FROM products p
		WHERE p.deleted_at IS NOT NULL
			AND (? = '' OR p.name ILIKE ? OR p.description ILIKE ?)
		ORDER BY p.product_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, product_errors.ErrFindTrashedProducts.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at,
			COUNT(*) OVER() as total_count
		FROM products p
		WHERE p.deleted_at IS NULL
			AND p.merchant_id = ?
			AND (? = '' OR p.name ILIKE ?)
			AND (? = 0 OR p.category_id = ?)
			AND (? = 0 OR p.price >= ?)
			AND (? = 0 OR p.price <= ?)
		ORDER BY p.product_id DESC
		LIMIT ? OFFSET ?
	`, req.MerchantID, req.Search, "%"+req.Search+"%",
		0, req.CategoryID,
		0, IntPtrToInt(req.MinPrice),
		0, IntPtrToInt(req.MaxPrice),
		req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, product_errors.ErrFindProductsByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*ProductResult, error) {
	offset := (req.Page - 1) * req.PageSize

	categoryID := 0
	if req.CategoryName != "" {
		cat, err := r.categoryRepo.FindByName(ctx, req.CategoryName)
		if err != nil {
			return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
		}
		if cat != nil {
			categoryID = int(cat.CategoryID)
		}
	}

	var results []*ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at,
			COUNT(*) OVER() as total_count
		FROM products p
		WHERE p.deleted_at IS NULL
			AND (? = 0 OR p.category_id = ?)
			AND (? = '' OR p.name ILIKE ?)
			AND (? = 0 OR p.price >= ?)
			AND (? = 0 OR p.price <= ?)
		ORDER BY p.product_id DESC
		LIMIT ? OFFSET ?
	`, categoryID, categoryID, req.Search, "%"+req.Search+"%",
		0, IntPtrToInt(req.MinPrice),
		0, IntPtrToInt(req.MaxPrice),
		req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, product_errors.ErrFindProductsByCategory.WithInternal(err)
	}
	return results, nil
}

func (r *productQueryRepository) FindByID(ctx context.Context, productID int) (*ProductResult, error) {
	var result ProductResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at
		FROM products p
		WHERE p.product_id = ? AND p.deleted_at IS NULL
	`, productID).Scan(&result).Error
	if err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	if result.ProductID == 0 {
		return nil, product_errors.ErrProductNotFound
	}
	return &result, nil
}

func (r *productQueryRepository) FindByIDTrashed(ctx context.Context, productID int) (*ProductResult, error) {
	var result ProductResult
	err := r.db.WithContext(ctx).Unscoped().Raw(`
		SELECT p.product_id, p.merchant_id, p.category_id, p.name, p.description,
			p.price, p.count_in_stock, p.brand, p.weight, p.rating, p.slug_product,
			p.image_product, p.created_at, p.updated_at, p.deleted_at
		FROM products p
		WHERE p.product_id = ?
	`, productID).Scan(&result).Error
	if err != nil {
		return nil, product_errors.ErrProductNotFound
	}
	return &result, nil
}
