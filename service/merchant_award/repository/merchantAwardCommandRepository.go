package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_award_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_award"
	"gorm.io/gorm"
)

type merchantAwardCommandRepository struct {
	db *gorm.DB
}

func NewMerchantAwardCommandRepository(db *gorm.DB) MerchantAwardCommandRepository {
	return &merchantAwardCommandRepository{db: db}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timePtr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, _ := time.Parse("2006-01-02", s)
	return &t
}

func (r *merchantAwardCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error) {
	now := time.Now()
	item := &models.MerchantCertificationsAndAward{
		MerchantID:     int32(req.MerchantID),
		Title:          req.Title,
		Description:    stringPtr(req.Description),
		IssuedBy:       stringPtr(req.IssuedBy),
		IssueDate:      timePtr(req.IssueDate),
		ExpiryDate:     timePtr(req.ExpiryDate),
		CertificateUrl: stringPtr(req.CertificateUrl),
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, merchant_award_errors.ErrCreateMerchantAward.WithInternal(err)
	}
	return item, nil
}

func (r *merchantAwardCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error) {
	var item models.MerchantCertificationsAndAward
	if err := r.db.WithContext(ctx).Where("merchant_certification_id = ? AND deleted_at IS NULL", *req.MerchantCertificationID).First(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrMerchantAwardNotFound
	}
	item.Title = req.Title
	item.Description = stringPtr(req.Description)
	item.IssuedBy = stringPtr(req.IssuedBy)
	item.IssueDate = timePtr(req.IssueDate)
	item.ExpiryDate = timePtr(req.ExpiryDate)
	item.CertificateUrl = stringPtr(req.CertificateUrl)
	now := time.Now()
	item.UpdatedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrUpdateMerchantAward.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantAwardCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error) {
	var item models.MerchantCertificationsAndAward
	if err := r.db.WithContext(ctx).Where("merchant_certification_id = ? AND deleted_at IS NULL", id).First(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrMerchantAwardNotFound
	}
	now := time.Now()
	item.DeletedAt = &now
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrTrashedMerchantAward.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantAwardCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error) {
	var item models.MerchantCertificationsAndAward
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_certification_id = ? AND deleted_at IS NOT NULL", id).First(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrMerchantAwardNotFound
	}
	item.DeletedAt = nil
	if err := r.db.WithContext(ctx).Unscoped().Save(&item).Error; err != nil {
		return nil, merchant_award_errors.ErrRestoreMerchantAward.WithInternal(err)
	}
	return &item, nil
}

func (r *merchantAwardCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_certification_id = ?", id).Delete(&models.MerchantCertificationsAndAward{})
	if result.Error != nil {
		return false, merchant_award_errors.ErrDeleteMerchantAwardPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_award_errors.ErrMerchantAwardNotFound
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantCertificationsAndAward{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_award_errors.ErrRestoreAllMerchantAwards.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantCertificationsAndAward{})
	if result.Error != nil {
		return false, merchant_award_errors.ErrDeleteAllMerchantAwardsPermanent.WithInternal(result.Error)
	}
	return true, nil
}
