package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_policy_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_policy_errors"
	"gorm.io/gorm"
)

type merchantPolicyQueryRepository struct {
	db *gorm.DB
}

func NewMerchantPolicyQueryRepository(db *gorm.DB) MerchantPoliciesQueryRepository {
	return &merchantPolicyQueryRepository{db: db}
}

func (r *merchantPolicyQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantPolicyResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mp.merchant_policy_id, mp.merchant_id, mp.policy_type, mp.title, mp.description,
			COALESCE(m.name, '') AS merchant_name,
			mp.created_at, mp.updated_at, mp.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_policies mp
		LEFT JOIN merchants m ON m.merchant_id = mp.merchant_id
		WHERE mp.deleted_at IS NULL
			AND (? = '' OR m.name ILIKE ? OR mp.policy_type ILIKE ?)
		ORDER BY mp.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_policy_errors.ErrFindAllMerchantPolicies.WithInternal(err)
	}
	return results, nil
}

func (r *merchantPolicyQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantPolicyResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mp.merchant_policy_id, mp.merchant_id, mp.policy_type, mp.title, mp.description,
			COALESCE(m.name, '') AS merchant_name,
			mp.created_at, mp.updated_at, mp.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_policies mp
		LEFT JOIN merchants m ON m.merchant_id = mp.merchant_id
		WHERE mp.deleted_at IS NULL
			AND (? = '' OR m.name ILIKE ? OR mp.policy_type ILIKE ?)
		ORDER BY mp.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_policy_errors.ErrFindActiveMerchantPolicies.WithInternal(err)
	}
	return results, nil
}

func (r *merchantPolicyQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantPolicyResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantPolicyResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mp.merchant_policy_id, mp.merchant_id, mp.policy_type, mp.title, mp.description,
			COALESCE(m.name, '') AS merchant_name,
			mp.created_at, mp.updated_at, mp.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_policies mp
		LEFT JOIN merchants m ON m.merchant_id = mp.merchant_id
		WHERE mp.deleted_at IS NOT NULL
			AND (? = '' OR m.name ILIKE ? OR mp.policy_type ILIKE ?)
		ORDER BY mp.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_policy_errors.ErrFindTrashedMerchantPolicies.WithInternal(err)
	}
	return results, nil
}

func (r *merchantPolicyQueryRepository) FindByID(ctx context.Context, id int) (*MerchantPolicyResult, error) {
	var result MerchantPolicyResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mp.merchant_policy_id, mp.merchant_id, mp.policy_type, mp.title, mp.description,
			COALESCE(m.name, '') AS merchant_name,
			mp.created_at, mp.updated_at, mp.deleted_at, 0 AS total_count
		FROM merchant_policies mp
		LEFT JOIN merchants m ON m.merchant_id = mp.merchant_id
		WHERE mp.merchant_policy_id = ? AND mp.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, merchant_policy_errors.ErrMerchantPolicyNotFound.WithInternal(err)
	}
	if result.MerchantPolicyID == 0 {
		return nil, merchant_policy_errors.ErrMerchantPolicyNotFound
	}
	return &result, nil
}

