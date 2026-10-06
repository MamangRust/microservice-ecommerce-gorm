package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_award_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant_award"
	"gorm.io/gorm"
)

type merchantAwardQueryRepository struct {
	db *gorm.DB
}

func NewMerchantAwardQueryRepository(db *gorm.DB) MerchantAwardQueryRepository {
	return &merchantAwardQueryRepository{db: db}
}

func (r *merchantAwardQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantCertResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mca.merchant_certification_id, mca.merchant_id, mca.title, mca.description,
			mca.issued_by, mca.issue_date, mca.expiry_date, mca.certificate_url,
			mca.created_at, mca.updated_at, mca.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_certifications_and_awards mca
		WHERE mca.deleted_at IS NULL
			AND (? = '' OR mca.title ILIKE ?)
		ORDER BY mca.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_award_errors.ErrFindAllMerchantAwards.WithInternal(err)
	}
	return results, nil
}

func (r *merchantAwardQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantCertResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mca.merchant_certification_id, mca.merchant_id, mca.title, mca.description,
			mca.issued_by, mca.issue_date, mca.expiry_date, mca.certificate_url,
			mca.created_at, mca.updated_at, mca.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_certifications_and_awards mca
		WHERE mca.deleted_at IS NULL
			AND (? = '' OR mca.title ILIKE ?)
		ORDER BY mca.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_award_errors.ErrFindByActiveMerchantAwards.WithInternal(err)
	}
	return results, nil
}

func (r *merchantAwardQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantCertResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mca.merchant_certification_id, mca.merchant_id, mca.title, mca.description,
			mca.issued_by, mca.issue_date, mca.expiry_date, mca.certificate_url,
			mca.created_at, mca.updated_at, mca.deleted_at,
			COUNT(*) OVER() AS total_count
		FROM merchant_certifications_and_awards mca
		WHERE mca.deleted_at IS NOT NULL
			AND (? = '' OR mca.title ILIKE ?)
		ORDER BY mca.created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%"+req.Search+"%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_award_errors.ErrFindByTrashedMerchantAwards.WithInternal(err)
	}
	return results, nil
}

func (r *merchantAwardQueryRepository) FindByID(ctx context.Context, id int) (*MerchantCertResult, error) {
	var result MerchantCertResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT mca.merchant_certification_id, mca.merchant_id, mca.title, mca.description,
			mca.issued_by, mca.issue_date, mca.expiry_date, mca.certificate_url,
			mca.created_at, mca.updated_at, mca.deleted_at, 0 AS total_count
		FROM merchant_certifications_and_awards mca
		WHERE mca.merchant_certification_id = ? AND mca.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, merchant_award_errors.ErrMerchantAwardNotFound.WithInternal(err)
	}
	if result.MerchantCertificationID == 0 {
		return nil, merchant_award_errors.ErrMerchantAwardNotFound
	}
	return &result, nil
}
