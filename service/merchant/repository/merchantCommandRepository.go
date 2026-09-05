package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shared_errors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"gorm.io/gorm"
)

type merchantCommandRepository struct {
	db *gorm.DB
}

func NewMerchantCommandRepository(db *gorm.DB) *merchantCommandRepository {
	return &merchantCommandRepository{db: db}
}

func (r *merchantCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	merchant := &models.Merchant{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  stringPtr(request.Description),
		Address:      stringPtr(request.Address),
		ContactEmail: stringPtr(request.ContactEmail),
		ContactPhone: stringPtr(request.ContactPhone),
	}

	if err := r.db.WithContext(ctx).Create(merchant).Error; err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantCommandRepository) CreateInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	merchant := &models.Merchant{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Status:       "active",
		Description:  stringPtr(request.Description),
		Address:      stringPtr(request.Address),
		ContactEmail: stringPtr(request.ContactEmail),
		ContactPhone: stringPtr(request.ContactPhone),
	}

	if err := tx.WithContext(ctx).Create(merchant).Error; err != nil {
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).First(&merchant, *request.MerchantID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	updates := map[string]interface{}{
		"name":          request.Name,
		"description":   stringPtr(request.Description),
		"address":       stringPtr(request.Address),
		"contact_email": stringPtr(request.ContactEmail),
		"contact_phone": stringPtr(request.ContactPhone),
		"status":        request.Status,
	}

	if err := r.db.WithContext(ctx).Model(&merchant).Updates(updates).Error; err != nil {
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&merchant, *request.MerchantID)
	return &merchant, nil
}

func (r *merchantCommandRepository) Trash(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).First(&merchant, merchantID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	if err := r.db.WithContext(ctx).Model(&merchant).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) Restore(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_id = ? AND deleted_at IS NOT NULL", merchantID).First(&merchant).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&merchant).Update("deleted_at", nil).Error; err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}
	r.db.WithContext(ctx).Unscoped().First(&merchant, merchantID)
	return &merchant, nil
}

func (r *merchantCommandRepository) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_id = ?", merchantID).Delete(&models.Merchant{})
	if result.Error != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete merchant while related records exist").WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_errors.ErrMerchantNotFound
	}
	return true, nil
}

func (r *merchantCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Merchant{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_errors.ErrRestoreAllMerchants.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Merchant{})
	if result.Error != nil {
		return false, shared_errors.NewConflictError("cannot permanently delete merchants while related records exist").WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).First(&merchant, *request.MerchantID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	if err := r.db.WithContext(ctx).Model(&merchant).Update("status", request.Status).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&merchant, *request.MerchantID)
	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := tx.WithContext(ctx).First(&merchant, *request.MerchantID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	if err := tx.WithContext(ctx).Model(&merchant).Update("status", request.Status).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	tx.WithContext(ctx).First(&merchant, *request.MerchantID)
	return &merchant, nil
}
