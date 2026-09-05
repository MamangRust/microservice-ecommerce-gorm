package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/order_errors"
	"gorm.io/gorm"
)

type orderCommandRepository struct {
	db *gorm.DB
}

func NewOrderCommandRepository(db *gorm.DB) OrderCommandRepository {
	return &orderCommandRepository{db: db}
}

func (r *orderCommandRepository) Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error) {
	order := &models.Order{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}
	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}
	return order, nil
}

func (r *orderCommandRepository) Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Where("order_id = ? AND deleted_at IS NULL", request.OrderID).First(&order).Error
	if err != nil {
		return nil, order_errors.ErrOrderNotFound
	}
	order.TotalPrice = int32(request.TotalPrice)
	err = r.db.WithContext(ctx).Save(&order).Error
	if err != nil {
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) Trash(ctx context.Context, order_id int) (*models.Order, error) {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("order_id = ? AND deleted_at IS NULL", order_id).
		Update("deleted_at", now)
	if result.Error != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, order_errors.ErrOrderNotFound
	}
	var order models.Order
	if err := r.db.WithContext(ctx).Where("order_id = ?", order_id).First(&order).Error; err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) Restore(ctx context.Context, order_id int) (*models.Order, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Order{}).
		Where("order_id = ? AND deleted_at IS NOT NULL", order_id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, order_errors.ErrOrderNotFound
	}
	var order models.Order
	if err := r.db.WithContext(ctx).Where("order_id = ?", order_id).First(&order).Error; err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashedByID(ctx context.Context, order_id int) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Unscoped().Where("order_id = ? AND deleted_at IS NOT NULL", order_id).First(&order).Error
	if err != nil {
		return nil, order_errors.ErrOrderNotFound
	}
	return &order, nil
}

func (r *orderCommandRepository) FindTrashed(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Find(&orders).Error
	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}
	return orders, nil
}

func (r *orderCommandRepository) DeletePermanent(ctx context.Context, order_id int) (bool, error) {
	// Delete child records first to avoid FK constraint
	r.db.WithContext(ctx).Unscoped().Where("order_id = ?", order_id).Delete(&models.OrderStockReservation{})
	
	result := r.db.WithContext(ctx).Unscoped().Where("order_id = ? AND deleted_at IS NOT NULL", order_id).Delete(&models.Order{})
	if result.Error != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, order_errors.ErrOrderNotFound
	}
	return true, nil
}

func (r *orderCommandRepository) DeletePermanentWithChildren(ctx context.Context, order_id int) (bool, error) {
	return r.DeletePermanent(ctx, order_id)
}

func (r *orderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().Model(&models.Order{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil).Error
	if err != nil {
		return false, order_errors.ErrRestoreAllOrder.WithInternal(err)
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Order{}).Error
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
