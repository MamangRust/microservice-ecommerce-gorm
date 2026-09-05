package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stockReservationRepository struct {
	db *gorm.DB
}

func NewStockReservationRepository(db *gorm.DB) StockReservationRepository {
	return &stockReservationRepository{db: db}
}

func (r *stockReservationRepository) GetByOrder(ctx context.Context, orderID int) ([]*models.OrderStockReservation, error) {
	var reservations []*models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&reservations).Error
	return reservations, err
}

func (r *stockReservationRepository) Upsert(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error) {
	var reservation models.OrderStockReservation
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"quantity", "status", "updated_at"}),
	}).Create(&models.OrderStockReservation{
		OrderID: int32(orderID), ProductID: int32(productID), Quantity: int32(quantity), Status: "reserved",
	}).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&reservation).Error
	return &reservation, err
}

func (r *stockReservationRepository) UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error) {
	var reservation models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&reservation).Error
	if err != nil { return nil, err }
	reservation.Quantity = int32(quantity)
	err = r.db.WithContext(ctx).Save(&reservation).Error
	return &reservation, err
}

func (r *stockReservationRepository) Release(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error) {
	result := r.db.WithContext(ctx).Model(&models.OrderStockReservation{}).
		Where("order_id = ? AND product_id = ? AND status = ?", orderID, productID, "reserved").
		Update("status", "released")
	if result.Error != nil { return nil, result.Error }
	if result.RowsAffected == 0 { return nil, gorm.ErrRecordNotFound }
	var reservation models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&reservation).Error
	return &reservation, err
}

func (r *stockReservationRepository) Reserve(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error) {
	result := r.db.WithContext(ctx).Model(&models.OrderStockReservation{}).
		Where("order_id = ? AND product_id = ? AND status = ?", orderID, productID, "released").
		Update("status", "reserved")
	if result.Error != nil { return nil, result.Error }
	if result.RowsAffected == 0 { return nil, gorm.ErrRecordNotFound }
	var reservation models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).First(&reservation).Error
	return &reservation, err
}

func (r *stockReservationRepository) GetReservedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	var reservations []*models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("status = ? AND order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL)", "reserved").Find(&reservations).Error
	return reservations, err
}

func (r *stockReservationRepository) GetReleasedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	var reservations []*models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("status = ? AND order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL)", "released").Find(&reservations).Error
	return reservations, err
}

func (r *stockReservationRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&models.OrderStockReservation{}).Error
}

func (r *stockReservationRepository) DeleteByOrderProduct(ctx context.Context, orderID, productID int) error {
	return r.db.WithContext(ctx).Where("order_id = ? AND product_id = ?", orderID, productID).Delete(&models.OrderStockReservation{}).Error
}

func (r *stockReservationRepository) DeleteAllForTrashedOrders(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL)").Delete(&models.OrderStockReservation{}).Error
}

func (r *stockReservationRepository) GetReleasedForActiveOrders(ctx context.Context) ([]*models.OrderStockReservation, error) {
	var reservations []*models.OrderStockReservation
	err := r.db.WithContext(ctx).Where("status = ? AND order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NULL)", "released").Find(&reservations).Error
	return reservations, err
}

func (r *stockReservationRepository) DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("status = ? AND updated_at < ? AND order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL)", "released", cutoff).Delete(&models.OrderStockReservation{})
	return result.RowsAffected, result.Error
}

func (r *stockReservationRepository) DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&models.ProductStockAdjustment{})
	return result.RowsAffected, result.Error
}
