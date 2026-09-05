package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/slider_errors"
	"gorm.io/gorm"
)

type sliderQueryRepository struct {
	db *gorm.DB
}

func NewSliderQueryRepository(db *gorm.DB) *sliderQueryRepository {
	return &sliderQueryRepository{db: db}
}

func (r *sliderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*SliderResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.slider_id, s.name, s.image,
			COALESCE(TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(s.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(s.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM sliders s
		WHERE s.deleted_at IS NULL
			AND (? = '' OR s.name ILIKE ?)
		ORDER BY s.slider_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, slider_errors.ErrFindAllSliders
	}
	return results, nil
}

func (r *sliderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*SliderResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.slider_id, s.name, s.image,
			COALESCE(TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(s.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(s.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM sliders s
		WHERE s.deleted_at IS NULL
			AND (? = '' OR s.name ILIKE ?)
		ORDER BY s.slider_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, slider_errors.ErrFindActiveSliders
	}
	return results, nil
}

func (r *sliderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*SliderResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*SliderResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.slider_id, s.name, s.image,
			COALESCE(TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(s.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(s.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM sliders s
		WHERE s.deleted_at IS NOT NULL
			AND (? = '' OR s.name ILIKE ?)
		ORDER BY s.slider_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, slider_errors.ErrFindTrashedSliders
	}
	return results, nil
}

func (r *sliderQueryRepository) FindByID(ctx context.Context, sliderID int) (*SliderResult, error) {
	var result SliderResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT s.slider_id, s.name, s.image,
			COALESCE(TO_CHAR(s.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(s.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(s.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at
		FROM sliders s
		WHERE s.slider_id = ? AND s.deleted_at IS NULL
	`, sliderID).Scan(&result).Error
	if err != nil {
		return nil, slider_errors.ErrFindSliderByID
	}
	if result.SliderID == 0 {
		return nil, slider_errors.ErrFindSliderByID
	}
	return &result, nil
}
