package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"gorm.io/gorm"
)

type shippingAddressCommandRepository struct {
	db *gorm.DB
}

func NewShippingAddressCommandRepository(db *gorm.DB) *shippingAddressCommandRepository {
	return &shippingAddressCommandRepository{db: db}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *shippingAddressCommandRepository) Create(ctx context.Context, request *requests.CreateShippingAddressRequest) (*models.ShippingAddress, error) {
	var orderID int32
	if request.OrderID != nil {
		orderID = int32(*request.OrderID)
	}

	addr := &models.ShippingAddress{
		OrderID:        orderID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   float64(request.ShippingCost),
		CreatedAt:      timePtr(time.Now()),
		UpdatedAt:      timePtr(time.Now()),
	}

	if err := r.db.WithContext(ctx).Create(addr).Error; err != nil {
		return nil, shippingaddress_errors.ErrCreateShippingAddress
	}
	return addr, nil
}

func (r *shippingAddressCommandRepository) Update(ctx context.Context, request *requests.UpdateShippingAddressRequest) (*models.ShippingAddress, error) {
	var addr models.ShippingAddress
	if err := r.db.WithContext(ctx).First(&addr, *request.ShippingID).Error; err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}

	addr.Alamat = request.Alamat
	addr.Provinsi = request.Provinsi
	addr.Kota = request.Kota
	addr.Negara = request.Negara
	addr.Courier = request.Courier
	addr.ShippingMethod = request.ShippingMethod
	addr.ShippingCost = float64(request.ShippingCost)
	addr.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&addr).Error; err != nil {
		return nil, shippingaddress_errors.ErrUpdateShippingAddress
	}
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Trash(ctx context.Context, shippingID int) (*models.ShippingAddress, error) {
	var addr models.ShippingAddress
	if err := r.db.WithContext(ctx).First(&addr, shippingID).Error; err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	now := timePtr(time.Now())
	if err := r.db.WithContext(ctx).Model(&addr).Update("deleted_at", now).Error; err != nil {
		return nil, shippingaddress_errors.ErrTrashShippingAddress
	}
	addr.DeletedAt = now
	return &addr, nil
}

func (r *shippingAddressCommandRepository) Restore(ctx context.Context, shippingID int) (*models.ShippingAddress, error) {
	var addr models.ShippingAddress
	if err := r.db.WithContext(ctx).Unscoped().Where("shipping_address_id = ?", shippingID).First(&addr).Error; err != nil {
		return nil, shippingaddress_errors.ErrShippingAddressNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&addr).Update("deleted_at", nil).Error; err != nil {
		return nil, shippingaddress_errors.ErrRestoreShippingAddress
	}
	addr.DeletedAt = nil
	return &addr, nil
}

func (r *shippingAddressCommandRepository) DeletePermanent(ctx context.Context, shippingID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("shipping_address_id = ?", shippingID).Delete(&models.ShippingAddress{})
	if result.Error != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	if result.RowsAffected == 0 {
		return false, shippingaddress_errors.ErrShippingAddressNotFound
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("order_id = ?", orderID).Delete(&models.ShippingAddress{})
	if result.Error != nil {
		return false, shippingaddress_errors.ErrDeleteShippingAddressPermanent
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Model(&models.ShippingAddress{}).Update("deleted_at", nil)
	if result.Error != nil {
		return false, shippingaddress_errors.ErrRestoreAllShippingAddresses
	}
	return true, nil
}

func (r *shippingAddressCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.ShippingAddress{})
	if result.Error != nil {
		return false, shippingaddress_errors.ErrDeleteAllPermanentShippingAddress
	}
	return true, nil
}
