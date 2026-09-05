package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"gorm.io/gorm"
)

type reviewDetailQueryRepository struct {
	db *gorm.DB
}

func NewReviewDetailQueryRepository(db *gorm.DB) *reviewDetailQueryRepository {
	return &reviewDetailQueryRepository{db: db}
}

func (r *reviewDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT rd.review_detail_id, rd.review_id, rd.type, rd.url, rd.caption,
			rd.created_at, rd.updated_at, rd.deleted_at,
			COUNT(*) OVER() as total_count
		FROM review_details rd
		WHERE rd.deleted_at IS NULL
			AND (? = '' OR rd.type ILIKE ? OR rd.caption ILIKE ?)
		ORDER BY rd.review_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_detail_errors.ErrFindAllReviewDetails.WithInternal(err)
	}
	return results, nil
}

func (r *reviewDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT rd.review_detail_id, rd.review_id, rd.type, rd.url, rd.caption,
			rd.created_at, rd.updated_at, rd.deleted_at,
			COUNT(*) OVER() as total_count
		FROM review_details rd
		WHERE rd.deleted_at IS NULL
			AND (? = '' OR rd.type ILIKE ? OR rd.caption ILIKE ?)
		ORDER BY rd.review_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_detail_errors.ErrFindActiveReviewDetails.WithInternal(err)
	}
	return results, nil
}

func (r *reviewDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewDetailResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT rd.review_detail_id, rd.review_id, rd.type, rd.url, rd.caption,
			rd.created_at, rd.updated_at, rd.deleted_at,
			COUNT(*) OVER() as total_count
		FROM review_details rd
		WHERE rd.deleted_at IS NOT NULL
			AND (? = '' OR rd.type ILIKE ? OR rd.caption ILIKE ?)
		ORDER BY rd.review_detail_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_detail_errors.ErrFindTrashedReviewDetails.WithInternal(err)
	}
	return results, nil
}

func (r *reviewDetailQueryRepository) FindByID(ctx context.Context, id int) (*ReviewDetailResult, error) {
	var result ReviewDetailResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT rd.review_detail_id, rd.review_id, rd.type, rd.url, rd.caption,
			rd.created_at, rd.updated_at, rd.deleted_at
		FROM review_details rd
		WHERE rd.review_detail_id = ?
	`, id).Scan(&result).Error
	if err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	if result.ReviewDetailID == 0 {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	return &result, nil
}
