package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review"
	"gorm.io/gorm"
)

type reviewCommandRepository struct {
	db *gorm.DB
}

func NewReviewCommandRepository(db *gorm.DB) *reviewCommandRepository {
	return &reviewCommandRepository{db: db}
}

func (r *reviewCommandRepository) Create(ctx context.Context, request *requests.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),

		Comment:   request.Comment,
		Rating:    int32(request.Rating),
	}
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return nil, review_errors.ErrCreateReview.WithInternal(err)
	}
	return review, nil
}

func (r *reviewCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewRequest) (*models.Review, error) {
	var review models.Review
	if err := r.db.WithContext(ctx).First(&review, *request.ReviewID).Error; err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	updates := map[string]interface{}{
		"name":    request.Name,
		"rating":  int32(request.Rating),
		"comment": request.Comment,
	}
	if err := r.db.WithContext(ctx).Model(&review).Updates(updates).Error; err != nil {
		return nil, review_errors.ErrUpdateReview.WithInternal(err)
	}
	r.db.WithContext(ctx).First(&review, *request.ReviewID)
	return &review, nil
}

func (r *reviewCommandRepository) Trash(ctx context.Context, reviewID int) (*models.Review, error) {
	var review models.Review
	if err := r.db.WithContext(ctx).First(&review, reviewID).Error; err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	if err := r.db.WithContext(ctx).Model(&review).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, review_errors.ErrTrashReview.WithInternal(err)
	}
	return &review, nil
}

func (r *reviewCommandRepository) Restore(ctx context.Context, reviewID int) (*models.Review, error) {
	var review models.Review
	if err := r.db.WithContext(ctx).Unscoped().Where("review_id = ? AND deleted_at IS NOT NULL", reviewID).First(&review).Error; err != nil {
		return nil, review_errors.ErrReviewNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&review).Update("deleted_at", nil).Error; err != nil {
		return nil, review_errors.ErrRestoreReview.WithInternal(err)
	}
	r.db.WithContext(ctx).Unscoped().First(&review, reviewID)
	return &review, nil
}

func (r *reviewCommandRepository) DeletePermanent(ctx context.Context, reviewID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("review_id = ?", reviewID).Delete(&models.Review{})
	if result.Error != nil {
		return false, review_errors.ErrDeleteReviewPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, review_errors.ErrReviewNotFound
	}
	return true, nil
}

func (r *reviewCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Review{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, review_errors.ErrRestoreAllReviews.WithInternal(result.Error)
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Review{})
	if result.Error != nil {
		return false, review_errors.ErrDeleteAllPermanentReview.WithInternal(result.Error)
	}
	return true, nil
}
