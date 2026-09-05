package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/order_errors"
	"gorm.io/gorm"
)

type orderQueryRepository struct {
	db *gorm.DB
}

func NewOrderQueryRepository(db *gorm.DB) OrderQueryRepository {
	return &orderQueryRepository{db: db}
}

func (r *orderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderResult
	query := `
		SELECT order_id, user_id, merchant_id, total_price, created_at, updated_at,
			COUNT(*) OVER() AS total_count
		FROM orders WHERE deleted_at IS NULL
			AND (? = '' OR CAST(order_id AS TEXT) ILIKE ?)
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`
	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}
	return results, nil
}

func (r *orderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderResult
	query := `
		SELECT order_id, user_id, merchant_id, total_price, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM orders WHERE deleted_at IS NULL
			AND (? = '' OR CAST(order_id AS TEXT) ILIKE ?)
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`
	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, order_errors.ErrFindByActive.WithInternal(err)
	}
	return results, nil
}

func (r *orderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderResult
	query := `
		SELECT order_id, user_id, merchant_id, total_price, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM orders WHERE deleted_at IS NOT NULL
			AND (? = '' OR CAST(order_id AS TEXT) ILIKE ?)
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`
	err := r.db.WithContext(ctx).Raw(query, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	return results, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderResult
	query := `
		SELECT order_id, user_id, merchant_id, total_price, created_at, updated_at,
			COUNT(*) OVER() AS total_count
		FROM orders WHERE deleted_at IS NULL AND merchant_id = ?
			AND (? = '' OR CAST(order_id AS TEXT) ILIKE ?)
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`
	err := r.db.WithContext(ctx).Raw(query, req.MerchantID, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, order_errors.ErrFindByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *orderQueryRepository) FindByID(ctx context.Context, id int) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Where("order_id = ? AND deleted_at IS NULL", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, order_errors.ErrOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrFindById.WithInternal(err)
	}
	return &order, nil
}
