package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"gorm.io/gorm"
)

type merchantDocumentQueryRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentQueryRepository(db *gorm.DB) *merchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{db: db}
}

func (r *merchantDocumentQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDocumentResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.document_id, md.merchant_id, md.document_type, md.document_url,
			md.status, md.note, md.uploaded_at, md.created_at, md.updated_at, md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_documents md
		WHERE md.deleted_at IS NULL
			AND (? = '' OR md.document_type ILIKE ? OR md.status ILIKE ?)
		ORDER BY md.document_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	var totalCount int
	if len(results) > 0 {
		totalCount = int(results[0].TotalCount)
	}
	return results, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindByID(ctx context.Context, id int) (*MerchantDocumentResult, error) {
	var result MerchantDocumentResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.document_id, md.merchant_id, md.document_type, md.document_url,
			md.status, md.note, md.uploaded_at, md.created_at, md.updated_at, md.deleted_at
		FROM merchant_documents md
		WHERE md.document_id = ?
	`, id).Scan(&result).Error
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}
	return &result, nil
}

func (r *merchantDocumentQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDocumentResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.document_id, md.merchant_id, md.document_type, md.document_url,
			md.status, md.note, md.uploaded_at, md.created_at, md.updated_at, md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_documents md
		WHERE md.deleted_at IS NULL
			AND (? = '' OR md.document_type ILIKE ? OR md.status ILIKE ?)
		ORDER BY md.document_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	var totalCount int
	if len(results) > 0 {
		totalCount = int(results[0].TotalCount)
	}
	return results, &totalCount, nil
}

func (r *merchantDocumentQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*MerchantDocumentResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			md.document_id, md.merchant_id, md.document_type, md.document_url,
			md.status, md.note, md.uploaded_at, md.created_at, md.updated_at, md.deleted_at,
			COUNT(*) OVER() as total_count
		FROM merchant_documents md
		WHERE md.deleted_at IS NOT NULL
			AND (? = '' OR md.document_type ILIKE ? OR md.status ILIKE ?)
		ORDER BY md.document_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	var totalCount int
	if len(results) > 0 {
		totalCount = int(results[0].TotalCount)
	}
	return results, &totalCount, nil
}
