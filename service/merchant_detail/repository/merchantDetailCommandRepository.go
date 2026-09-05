package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_detail"
	"gorm.io/gorm"
)

type merchantDetailCommandRepository struct {
	db *gorm.DB
}

func NewMerchantDetailCommandRepository(db *gorm.DB) *merchantDetailCommandRepository {
	return &merchantDetailCommandRepository{db: db}
}

func (r *merchantDetailCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*models.MerchantDetail, error) {
	detail := &models.MerchantDetail{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      &request.DisplayName,
		CoverImageUrl:    &request.CoverImageUrl,
		LogoUrl:          &request.LogoUrl,
		ShortDescription: &request.ShortDescription,
		WebsiteUrl:       &request.WebsiteUrl,
	}

	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail.WithInternal(err)
	}

	return detail, nil
}

func (r *merchantDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*models.MerchantDetail, error) {
	var detail models.MerchantDetail
	if err := r.db.WithContext(ctx).First(&detail, *request.MerchantDetailID).Error; err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}

	updates := map[string]interface{}{
		"display_name":      &request.DisplayName,
		"cover_image_url":   &request.CoverImageUrl,
		"logo_url":          &request.LogoUrl,
		"short_description": &request.ShortDescription,
		"website_url":       &request.WebsiteUrl,
	}

	if err := r.db.WithContext(ctx).Model(&detail).Updates(updates).Error; err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&detail, *request.MerchantDetailID)
	return &detail, nil
}

func (r *merchantDetailCommandRepository) Trash(ctx context.Context, merchantDetailID int) (*models.MerchantDetail, error) {
	var detail models.MerchantDetail
	if err := r.db.WithContext(ctx).First(&detail, merchantDetailID).Error; err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	if err := r.db.WithContext(ctx).Model(&detail).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, merchantdetail_errors.ErrTrashMerchantDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *merchantDetailCommandRepository) Restore(ctx context.Context, merchantDetailID int) (*models.MerchantDetail, error) {
	var detail models.MerchantDetail
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_detail_id = ? AND deleted_at IS NOT NULL", merchantDetailID).First(&detail).Error; err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&detail).Update("deleted_at", nil).Error; err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail.WithInternal(err)
	}
	r.db.WithContext(ctx).Unscoped().First(&detail, merchantDetailID)
	return &detail, nil
}

func (r *merchantDetailCommandRepository) DeletePermanent(ctx context.Context, merchantDetailID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_detail_id = ?", merchantDetailID).Delete(&models.MerchantDetail{})
	if result.Error != nil {
		return false, merchantdetail_errors.ErrDeletePermanentMerchantDetail.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantDetail{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantDetail{})
	if result.Error != nil {
		return false, merchantdetail_errors.ErrDeleteAllPermanentMerchantDetails.WithInternal(result.Error)
	}
	return true, nil
}
