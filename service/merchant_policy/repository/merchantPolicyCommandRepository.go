package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_policy_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_policy_errors"
	"gorm.io/gorm"
)

type merchantPolicyCommandRepository struct {
	db *gorm.DB
}

func NewMerchantPolicyCommandRepository(db *gorm.DB) MerchantPoliciesCommandRepository {
	return &merchantPolicyCommandRepository{db: db}
}

func (r *merchantPolicyCommandRepository) Create(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*models.MerchantPolicy, error) {
	now := time.Now()
	policy := &models.MerchantPolicy{
		MerchantID:  int32(req.MerchantID),
		PolicyType:  req.PolicyType,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	if err := r.db.WithContext(ctx).Create(policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrCreateMerchantPolicy.WithInternal(err)
	}
	return policy, nil
}

func (r *merchantPolicyCommandRepository) Update(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*models.MerchantPolicy, error) {
	var policy models.MerchantPolicy
	if err := r.db.WithContext(ctx).Where("merchant_policy_id = ? AND deleted_at IS NULL", *req.MerchantPolicyID).First(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrMerchantPolicyNotFound
	}
	policy.PolicyType = req.PolicyType
	policy.Title = req.Title
	policy.Description = req.Description
	now := time.Now()
	policy.UpdatedAt = &now
	if err := r.db.WithContext(ctx).Save(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrUpdateMerchantPolicy.WithInternal(err)
	}
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) Trash(ctx context.Context, id int) (*models.MerchantPolicy, error) {
	var policy models.MerchantPolicy
	if err := r.db.WithContext(ctx).Where("merchant_policy_id = ? AND deleted_at IS NULL", id).First(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrMerchantPolicyNotFound
	}
	now := time.Now()
	policy.DeletedAt = &now
	if err := r.db.WithContext(ctx).Save(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrTrashedMerchantPolicy.WithInternal(err)
	}
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) Restore(ctx context.Context, id int) (*models.MerchantPolicy, error) {
	var policy models.MerchantPolicy
	if err := r.db.WithContext(ctx).Unscoped().Where("merchant_policy_id = ? AND deleted_at IS NOT NULL", id).First(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrMerchantPolicyNotFound
	}
	policy.DeletedAt = nil
	if err := r.db.WithContext(ctx).Unscoped().Save(&policy).Error; err != nil {
		return nil, merchant_policy_errors.ErrRestoreMerchantPolicy.WithInternal(err)
	}
	return &policy, nil
}

func (r *merchantPolicyCommandRepository) DeletePermanent(ctx context.Context, id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("merchant_policy_id = ?", id).Delete(&models.MerchantPolicy{})
	if result.Error != nil {
		return false, merchant_policy_errors.ErrDeleteMerchantPolicyPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_policy_errors.ErrMerchantPolicyNotFound
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantPolicy{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_policy_errors.ErrRestoreAllMerchantPolicies.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantPolicy{})
	if result.Error != nil {
		return false, merchant_policy_errors.ErrDeleteAllMerchantPoliciesPermanent.WithInternal(result.Error)
	}
	return true, nil
}
