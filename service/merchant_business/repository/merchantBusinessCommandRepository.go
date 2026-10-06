package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_business_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_business"
	"gorm.io/gorm"
)

type merchantBusinessCommandRepository struct {
	db *gorm.DB
}

func NewMerchantBusinessCommandRepository(db *gorm.DB) MerchantBusinessCommandRepository {
	return &merchantBusinessCommandRepository{db: db}
}

func int32Ptr(v int) *int32 {
	i := int32(v)
	return &i
}

func (r *merchantBusinessCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error) {
	now := time.Now()
	item := &models.MerchantBusinessInformation{
		MerchantID:        int32(req.MerchantID),
		BusinessType:      &req.BusinessType,
		TaxID:             &req.TaxID,
		EstablishedYear:   int32Ptr(req.EstablishedYear),
		NumberOfEmployees: int32Ptr(req.NumberOfEmployees),
		WebsiteUrl:        &req.WebsiteUrl,
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, merchant_business_errors.ErrCreateMerchantBusiness.WithInternal(err)
	}
	return item, nil
}

func (r *merchantBusinessCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*models.MerchantBusinessInformation, error) {
	var item models.MerchantBusinessInformation
	if err := r.db.WithContext(ctx).Where("merchant_business_info_id = ? AND deleted_at IS NULL", *req.MerchantBusinessInfoID).First(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound
	}
	item.BusinessType = &req.BusinessType
	item.TaxID = &req.TaxID
	item.EstablishedYear = int32Ptr(req.EstablishedYear)
	item.NumberOfEmployees = int32Ptr(req.NumberOfEmployees)
	item.WebsiteUrl = &req.WebsiteUrl
	now := time.Now()
	item.UpdatedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrUpdateMerchantBusiness.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantBusinessCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantBusinessInformation, error) {
	var item models.MerchantBusinessInformation
	if err := r.db.WithContext(ctx).Where("merchant_business_info_id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound
	}
	now := time.Now()
	item.DeletedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrTrashMerchantBusiness.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantBusinessCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantBusinessInformation, error) {
	var item models.MerchantBusinessInformation
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_business_info_id = ? AND deleted_at IS NOT NULL", id).First(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound
	}
	item.DeletedAt = nil
	if err := r.db.WithContext(ctx).Unscoped().Save(&item).Error; err != nil {
		return nil, merchant_business_errors.ErrRestoreMerchantBusiness.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantBusinessCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_business_info_id = ?", id).Delete(&models.MerchantBusinessInformation{})
	if result.Error != nil {
		return false, merchant_business_errors.ErrDeletePermanentMerchantBusiness.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_business_errors.ErrMerchantBusinessNotFound
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantBusinessInformation{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_business_errors.ErrRestoreAllMerchantBusinesses.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantBusinessInformation{})
	if result.Error != nil {
		return false, merchant_business_errors.ErrDeleteAllPermanentMerchantBusinesses.WithInternal(result.Error)
	}
	return true, nil
}
