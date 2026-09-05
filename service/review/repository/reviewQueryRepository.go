package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review"
	"gorm.io/gorm"
)

type reviewQueryRepository struct {
	db *gorm.DB
}

func NewReviewQueryRepository(db *gorm.DB) *reviewQueryRepository {
	return &reviewQueryRepository{db: db}
}

func (r *reviewQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at,
			COUNT(*) OVER() as total_count
		FROM reviews r
		WHERE r.deleted_at IS NULL
			AND (? = '' OR r.name ILIKE ? OR r.comment ILIKE ?)
		ORDER BY r.review_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_errors.ErrFindAllReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at,
			COUNT(*) OVER() as total_count
		FROM reviews r
		WHERE r.deleted_at IS NULL
			AND r.product_id = ?
			AND (? = 0 OR r.rating = ?)
		ORDER BY r.review_id DESC
		LIMIT ? OFFSET ?
	`, req.ProductID, req.Rating, req.Rating, req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_errors.ErrFindReviewsByProduct.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at,
			COUNT(*) OVER() as total_count
		FROM reviews r
		INNER JOIN products p ON p.product_id = r.product_id
		WHERE r.deleted_at IS NULL
			AND p.merchant_id = ?
			AND (? = 0 OR r.rating = ?)
		ORDER BY r.review_id DESC
		LIMIT ? OFFSET ?
	`, req.MerchantID, req.Rating, req.Rating, req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_errors.ErrFindReviewsByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at,
			COUNT(*) OVER() as total_count
		FROM reviews r
		WHERE r.deleted_at IS NULL
			AND (? = '' OR r.name ILIKE ? OR r.comment ILIKE ?)
		ORDER BY r.review_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_errors.ErrFindActiveReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*ReviewResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at,
			COUNT(*) OVER() as total_count
		FROM reviews r
		WHERE r.deleted_at IS NOT NULL
			AND (? = '' OR r.name ILIKE ? OR r.comment ILIKE ?)
		ORDER BY r.review_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, review_errors.ErrFindTrashedReviews.WithInternal(err)
	}
	return results, nil
}

func (r *reviewQueryRepository) FindByID(ctx context.Context, id int) (*ReviewResult, error) {
	var result ReviewResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT r.review_id, r.user_id, r.product_id, r.name, r.comment, r.rating,
			r.created_at, r.updated_at, r.deleted_at
		FROM reviews r
		WHERE r.review_id = ? AND r.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	if result.ReviewID == 0 {
		return nil, review_errors.ErrReviewNotFound
	}
	return &result, nil
}
