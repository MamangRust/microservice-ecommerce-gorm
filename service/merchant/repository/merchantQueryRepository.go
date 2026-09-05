package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"gorm.io/gorm"
)

type merchantQueryRepository struct {
	db *gorm.DB
}

func NewMerchantQueryRepository(db *gorm.DB) *merchantQueryRepository {
	return &merchantQueryRepository{db: db}
}

func fmtTime(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func (r *merchantQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			m.merchant_id, m.user_id, m.name, m.description, m.address,
			m.contact_email, m.contact_phone, m.status,
			m.created_at, m.updated_at, m.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchants m
		WHERE m.deleted_at IS NULL
			AND (? = '' OR m.name ILIKE ? OR m.description ILIKE ?)
		ORDER BY m.merchant_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_errors.ErrFindAllMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			m.merchant_id, m.user_id, m.name, m.description, m.address,
			m.contact_email, m.contact_phone, m.status,
			m.created_at, m.updated_at, m.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchants m
		WHERE m.deleted_at IS NULL
			AND (? = '' OR m.name ILIKE ? OR m.description ILIKE ?)
		ORDER BY m.merchant_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_errors.ErrFindActiveMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			m.merchant_id, m.user_id, m.name, m.description, m.address,
			m.contact_email, m.contact_phone, m.status,
			m.created_at, m.updated_at, m.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchants m
		WHERE m.deleted_at IS NOT NULL
			AND (? = '' OR m.name ILIKE ? OR m.description ILIKE ?)
		ORDER BY m.merchant_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchants.WithInternal(err)
	}
	return results, nil
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, userID int) (*MerchantResult, error) {
	var result MerchantResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			m.merchant_id, m.user_id, m.name, m.description, m.address,
			m.contact_email, m.contact_phone, m.status,
			m.created_at, m.updated_at, m.deleted_at
		FROM merchants m
		WHERE m.merchant_id = ? AND m.deleted_at IS NULL
	`, userID).Scan(&result).Error
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	if result.MerchantID == 0 {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	return &result, nil
}
