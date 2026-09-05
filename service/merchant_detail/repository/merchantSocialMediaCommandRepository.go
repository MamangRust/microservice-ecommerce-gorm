package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_social_link_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_social_link_errors"
	"gorm.io/gorm"
)

type merchantSocialLinkCommandRepository struct {
	db *gorm.DB
}

func NewMerchantSocialLinkCommandRepository(db *gorm.DB) *merchantSocialLinkCommandRepository {
	return &merchantSocialLinkCommandRepository{db: db}
}

func (r *merchantSocialLinkCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error) {
	link := &models.MerchantSocialMediaLink{
		MerchantDetailID: int32(*req.MerchantDetailID),
		Platform:         req.Platform,
		Url:              req.Url,
	}

	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return nil, merchant_social_link_errors.ErrCreateMerchantSocialLink.WithInternal(err)
	}

	return link, nil
}

func (r *merchantSocialLinkCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error) {
	var link models.MerchantSocialMediaLink
	if err := r.db.WithContext(ctx).First(&link, req.ID).Error; err != nil {
		return nil, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}

	updates := map[string]interface{}{
		"platform": req.Platform,
		"url":      req.Url,
	}

	if err := r.db.WithContext(ctx).Model(&link).Updates(updates).Error; err != nil {
		return nil, merchant_social_link_errors.ErrUpdateMerchantSocialLink.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&link, req.ID)
	return &link, nil
}

func (r *merchantSocialLinkCommandRepository) Trash(ctx context.Context, socialID int) (bool, error) {
	var link models.MerchantSocialMediaLink
	if err := r.db.WithContext(ctx).First(&link, socialID).Error; err != nil {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	if err := r.db.WithContext(ctx).Model(&link).Update("deleted_at", time.Now()).Error; err != nil {
		return false, merchant_social_link_errors.ErrTrashMerchantSocialLink.WithInternal(err)
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) Restore(ctx context.Context, socialID int) (bool, error) {
	var link models.MerchantSocialMediaLink
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_social_id = ? AND deleted_at IS NOT NULL", socialID).First(&link).Error; err != nil {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&link).Update("deleted_at", nil).Error; err != nil {
		return false, merchant_social_link_errors.ErrRestoreMerchantSocialLink.WithInternal(err)
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeletePermanent(ctx context.Context, socialID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_social_id = ?", socialID).Delete(&models.MerchantSocialMediaLink{})
	if result.Error != nil {
		return false, merchant_social_link_errors.ErrDeletePermanentMerchantSocialLink.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_social_link_errors.ErrMerchantSocialLinkNotFound
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantSocialMediaLink{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_social_link_errors.ErrRestoreAllMerchantSocialLinks.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantSocialLinkCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantSocialMediaLink{})
	if result.Error != nil {
		return false, merchant_social_link_errors.ErrDeleteAllPermanentMerchantSocialLinks.WithInternal(result.Error)
	}
	return true, nil
}
