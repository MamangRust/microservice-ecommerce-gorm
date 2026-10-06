package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_business_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_business"
	"gorm.io/gorm"
)

type merchantBusinessQueryRepository struct {
	db *gorm.DB
}

func NewMerchantBusinessQueryRepository(db *gorm.DB) MerchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{db: db}
}

func (r *merchantBusinessQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantBusinessResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mbi.merchant_business_info_id, mbi.merchant_id, mbi.business_type, mbi.tax_id,
			mbi.established_year, mbi.number_of_employees, mbi.website_url,
			mbi.created_at, mbi.updated_at, mbi.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_business_information mbi
		WHERE mbi.deleted_at IS NULL
			AND (? = '' OR mbi.business_type ILIKE ?)
		ORDER BY mbi.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound.WithInternal(err)
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantBusinessResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mbi.merchant_business_info_id, mbi.merchant_id, mbi.business_type, mbi.tax_id,
			mbi.established_year, mbi.number_of_employees, mbi.website_url,
			mbi.created_at, mbi.updated_at, mbi.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_business_information mbi
		WHERE mbi.deleted_at IS NULL
			AND (? = '' OR mbi.business_type ILIKE ?)
		ORDER BY mbi.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_business_errors.ErrFindActiveMerchantBusinesses.WithInternal(err)
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantBusinessResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantBusinessResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mbi.merchant_business_info_id, mbi.merchant_id, mbi.business_type, mbi.tax_id,
			mbi.established_year, mbi.number_of_employees, mbi.website_url,
			mbi.created_at, mbi.updated_at, mbi.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_business_information mbi
		WHERE mbi.deleted_at IS NOT NULL
			AND (? = '' OR mbi.business_type ILIKE ?)
		ORDER BY mbi.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_business_errors.ErrFindTrashedMerchantBusinesses.WithInternal(err)
	}
	return results, nil
}

func (r *merchantBusinessQueryRepository) FindByID(ctx context.Context, id int) (*MerchantBusinessResult, error) {
	var result MerchantBusinessResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mbi.merchant_business_info_id, mbi.merchant_id, mbi.business_type, mbi.tax_id,
			mbi.established_year, mbi.number_of_employees, mbi.website_url,
			mbi.created_at, mbi.updated_at, mbi.deleted_at, 0 AS total_count
		FROM merchant_business_information mbi
		WHERE mbi.merchant_business_info_id = ? AND mbi.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound.WithInternal(err)
	}
	if result.MerchantBusinessInfoID == 0 {
		return nil, merchant_business_errors.ErrMerchantBusinessNotFound
	}
	return &result, nil
}
