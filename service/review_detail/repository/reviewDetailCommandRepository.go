package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	review_detail_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/review_detail"
	"gorm.io/gorm"
)

type reviewDetailCommandRepository struct {
	db *gorm.DB
}

func NewReviewDetailCommandRepository(db *gorm.DB) *reviewDetailCommandRepository {
	return &reviewDetailCommandRepository{db: db}
}

func (r *reviewDetailCommandRepository) Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*models.ReviewDetail, error) {
	caption := request.Caption
	detail := &models.ReviewDetail{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  &caption,
	}
	if err := r.db.WithContext(ctx).Create(detail).Error; err != nil {
		return nil, review_detail_errors.ErrCreateReviewDetail.WithInternal(err)
	}
	return detail, nil
}

func (r *reviewDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*models.ReviewDetail, error) {
	var detail models.ReviewDetail
	if err := r.db.WithContext(ctx).First(&detail, *request.ReviewDetailID).Error; err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	caption := request.Caption
	updates := map[string]interface{}{
		"type":    request.Type,
		"url":     request.Url,
		"caption": &caption,
	}
	if err := r.db.WithContext(ctx).Model(&detail).Updates(updates).Error; err != nil {
		return nil, review_detail_errors.ErrUpdateReviewDetail.WithInternal(err)
	}
	r.db.WithContext(ctx).First(&detail, *request.ReviewDetailID)
	return &detail, nil
}

func (r *reviewDetailCommandRepository) Trash(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error) {
	var detail models.ReviewDetail
	if err := r.db.WithContext(ctx).First(&detail, reviewDetailID).Error; err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	if err := r.db.WithContext(ctx).Model(&detail).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, review_detail_errors.ErrTrashedReviewDetail.WithInternal(err)
	}
	return &detail, nil
}

func (r *reviewDetailCommandRepository) Restore(ctx context.Context, reviewDetailID int) (*models.ReviewDetail, error) {
	var detail models.ReviewDetail
	if err := r.db.WithContext(ctx).Unscoped().Where("review_detail_id = ? AND deleted_at IS NOT NULL", reviewDetailID).First(&detail).Error; err != nil {
		return nil, review_detail_errors.ErrReviewDetailNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&detail).Update("deleted_at", nil).Error; err != nil {
		return nil, review_detail_errors.ErrRestoreReviewDetail.WithInternal(err)
	}
	r.db.WithContext(ctx).Unscoped().First(&detail, reviewDetailID)
	return &detail, nil
}

func (r *reviewDetailCommandRepository) DeletePermanent(ctx context.Context, reviewDetailID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("review_detail_id = ?", reviewDetailID).Delete(&models.ReviewDetail{})
	if result.Error != nil {
		return false, review_detail_errors.ErrDeleteReviewDetailPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, review_detail_errors.ErrReviewDetailNotFound
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.ReviewDetail{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, review_detail_errors.ErrRestoreAllReviewDetails.WithInternal(result.Error)
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.ReviewDetail{})
	if result.Error != nil {
		return false, review_detail_errors.ErrDeleteAllReviewDetails.WithInternal(result.Error)
	}
	return true, nil
}
