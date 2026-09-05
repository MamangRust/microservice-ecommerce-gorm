package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/product_errors"
	"gorm.io/gorm"
)

type productCommandRepository struct {
	db *gorm.DB
}

func NewProductCommandRepository(db *gorm.DB) *productCommandRepository {
	return &productCommandRepository{db: db}
}

func (r *productCommandRepository) Create(ctx context.Context, request *requests.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       int32Ptr(request.Weight),
		SlugProduct:  request.SlugProduct,
		ImageProduct: &request.ImageProduct,
		CreatedAt:    timePtr(time.Now()),
		UpdatedAt:    timePtr(time.Now()),
	}

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, product_errors.ErrCreateProduct.WithInternal(err)
	}

	return product, nil
}

func (r *productCommandRepository) Update(ctx context.Context, request *requests.UpdateProductRequest) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, *request.ProductID).Error; err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	product.CategoryID = int32(request.CategoryID)
	product.Name = request.Name
	desc := request.Description
	product.Description = &desc
	product.Price = int32(request.Price)
	product.CountInStock = int32(request.CountInStock)
	brand := request.Brand
	product.Brand = &brand
	product.Weight = int32Ptr(request.Weight)
	product.SlugProduct = request.SlugProduct
	img := request.ImageProduct
	product.ImageProduct = &img
	product.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, product_errors.ErrUpdateProduct.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, product_id).Error; err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	product.CountInStock = int32(stock)
	product.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) AdjustProductStock(ctx context.Context, product_id int, delta int, operationID string) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, product_id).Error; err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	newStock := int(product.CountInStock) + delta
	if newStock < 0 {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	product.CountInStock = int32(newStock)
	product.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) Trash(ctx context.Context, product_id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, product_id).Error; err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	now := timePtr(time.Now())
	if err := r.db.WithContext(ctx).Model(&product).Update("deleted_at", now).Error; err != nil {
		return nil, product_errors.ErrTrashedProduct.WithInternal(err)
	}

	product.DeletedAt = now
	return &product, nil
}

func (r *productCommandRepository) Restore(ctx context.Context, product_id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Unscoped().Where("product_id = ?", product_id).First(&product).Error; err != nil {
		return nil, product_errors.ErrProductNotFound
	}

	if err := r.db.WithContext(ctx).Unscoped().Model(&product).Update("deleted_at", nil).Error; err != nil {
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}

	product.DeletedAt = nil
	return &product, nil
}

func (r *productCommandRepository) DeletePermanent(ctx context.Context, product_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("product_id = ?", product_id).Delete(&models.Product{})
	if result.Error != nil {
		if result.Error == gorm.ErrForeignKeyViolated {
			return false, shared_errors.NewConflictError("cannot permanently delete product while related records exist").WithInternal(result.Error)
		}
		return false, product_errors.ErrDeleteProductPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, product_errors.ErrProductNotFound
	}
	return true, nil
}

func (r *productCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Model(&models.Product{}).Update("deleted_at", nil)
	if result.Error != nil {
		return false, product_errors.ErrRestoreAllProducts.WithInternal(result.Error)
	}
	return true, nil
}

func (r *productCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Product{})
	if result.Error != nil {
		if result.Error == gorm.ErrForeignKeyViolated {
			return false, shared_errors.NewConflictError("cannot permanently delete products while related records exist").WithInternal(result.Error)
		}
		return false, product_errors.ErrDeleteAllProducts.WithInternal(result.Error)
	}
	return true, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
