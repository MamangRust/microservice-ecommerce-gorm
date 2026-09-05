package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/order_item_errors"
	"gorm.io/gorm"
)

type orderItemCommandRepository struct {
	db *gorm.DB
}

func NewOrderItemCommandRepository(db *gorm.DB) *orderItemCommandRepository {
	return &orderItemCommandRepository{db: db}
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	now := time.Now()
	item := &models.OrderItem{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err)
	}
	return item, nil
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	var item models.OrderItem
	if err := r.db.WithContext(ctx).Where("order_item_id = ? AND deleted_at IS NULL", req.OrderItemID).First(&item).Error; err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	item.Quantity = int32(req.Quantity)
	item.Price = int32(req.Price)
	now := time.Now()
	item.UpdatedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err)
	}
	return &item, nil
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, orderItemID int) (*models.OrderItem, error) {
	var item models.OrderItem
	if err := r.db.WithContext(ctx).Where("order_item_id = ? AND deleted_at IS NULL", orderItemID).First(&item).Error; err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	now := time.Now()
	item.DeletedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}
	return &item, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, orderItemID int) (*models.OrderItem, error) {
	var item models.OrderItem
	if err := r.db.WithContext(ctx).Unscoped().Where("order_item_id = ? AND deleted_at IS NOT NULL", orderItemID).First(&item).Error; err != nil {
		return nil, orderitem_errors.ErrOrderItemNotFound
	}
	item.DeletedAt = nil
	if err := r.db.WithContext(ctx).Unscoped().Save(&item).Error; err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}
	return &item, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, orderItemID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("order_item_id = ?", orderItemID).Delete(&models.OrderItem{})
	if result.Error != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	if result.RowsAffected == 0 {
		return false, orderitem_errors.ErrOrderItemNotFound
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("order_id = ?", orderID).Delete(&models.OrderItem{})
	if result.Error != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.OrderItem{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.OrderItem{})
	if result.Error != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) CalculateTotalPrice(ctx context.Context, orderID int) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.OrderItem{}).
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Select("COALESCE(SUM(price * quantity), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, orderitem_errors.ErrCalculateTotalPrice
	}
	return int(total), nil
}
