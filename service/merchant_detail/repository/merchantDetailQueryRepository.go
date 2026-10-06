package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_detail"
	"gorm.io/gorm"
)

type merchantDetailQueryRepository struct {
	db *gorm.DB
}

func NewMerchantDetailQueryRepository(db *gorm.DB) *merchantDetailQueryRepository {
	return &merchantDetailQueryRepository{db: db}
}

func (r *merchantDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.merchant_detail_id,
			md.merchant_id,
			md.display_name,
			md.cover_image_url,
			md.logo_url,
			md.short_description,
			md.website_url,
			md.created_at,
			md.updated_at,
			md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_details md
		WHERE md.deleted_at IS NULL
			AND (? = '' OR md.display_name ILIKE ? OR md.short_description ILIKE ?)
		ORDER BY md.merchant_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails.WithInternal(err)
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.merchant_detail_id,
			md.merchant_id,
			md.display_name,
			md.cover_image_url,
			md.logo_url,
			md.short_description,
			md.website_url,
			md.created_at,
			md.updated_at,
			md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_details md
		WHERE md.deleted_at IS NULL
			AND (? = '' OR md.display_name ILIKE ? OR md.short_description ILIKE ?)
		ORDER BY md.merchant_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchantdetail_errors.ErrFindActiveMerchantDetails.WithInternal(err)
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.merchant_detail_id,
			md.merchant_id,
			md.display_name,
			md.cover_image_url,
			md.logo_url,
			md.short_description,
			md.website_url,
			md.created_at,
			md.updated_at,
			md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_details md
		WHERE md.deleted_at IS NOT NULL
			AND (? = '' OR md.display_name ILIKE ? OR md.short_description ILIKE ?)
		ORDER BY md.merchant_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchantdetail_errors.ErrFindTrashedMerchantDetails.WithInternal(err)
	}
	return results, nil
}

func (r *merchantDetailQueryRepository) FindByID(ctx context.Context, userID int) (*MerchantDetailResult, error) {
	var result MerchantDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.merchant_detail_id,
			md.merchant_id,
			md.display_name,
			md.cover_image_url,
			md.logo_url,
			md.short_description,
			md.website_url,
			md.created_at,
			md.updated_at,
			md.deleted_at
		FROM merchant_details md
		WHERE md.merchant_detail_id = ? AND md.deleted_at IS NULL
	`, userID).Scan(&result).Error
	if err != nil {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound.WithInternal(err)
	}
	if result.MerchantDetailID == 0 {
		return nil, merchantdetail_errors.ErrMerchantDetailNotFound
	}
	return &result, nil
}
