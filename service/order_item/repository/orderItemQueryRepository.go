package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"gorm.io/gorm"
)

type orderItemQueryRepository struct {
	db *gorm.DB
}

func NewOrderItemQueryRepository(db *gorm.DB) *orderItemQueryRepository {
	return &orderItemQueryRepository{db: db}
}

func (r *orderItemQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderItemResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT order_item_id, order_id, product_id, quantity, price,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM order_items
		WHERE deleted_at IS NULL
			AND (? = '' OR CAST(product_id AS TEXT) LIKE '%' || ? || '%')
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems.WithInternal(err)
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderItemResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT order_item_id, order_id, product_id, quantity, price,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM order_items
		WHERE deleted_at IS NULL
			AND (? = '' OR CAST(product_id AS TEXT) LIKE '%' || ? || '%')
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, orderitem_errors.ErrFindByActive.WithInternal(err)
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*OrderItemResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT order_item_id, order_id, product_id, quantity, price,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM order_items
		WHERE deleted_at IS NOT NULL
			AND (? = '' OR CAST(product_id AS TEXT) LIKE '%' || ? || '%')
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, orderitem_errors.ErrFindByTrashed.WithInternal(err)
	}
	return results, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*OrderItemResult, error) {
	var results []*OrderItemResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT order_item_id, order_id, product_id, quantity, price,
			created_at, updated_at, deleted_at, 0 AS total_count
		FROM order_items
		WHERE order_id = ? AND deleted_at IS NULL
		ORDER BY order_item_id ASC
	`, order_id).Scan(&results).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, orderitem_errors.ErrFindOrderItemByOrder
		}
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}
	return results, nil
}
