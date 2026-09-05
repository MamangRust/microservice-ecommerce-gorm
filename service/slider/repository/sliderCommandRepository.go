package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/slider_errors"
	"gorm.io/gorm"
)

type sliderCommandRepository struct {
	db *gorm.DB
}

func NewSliderCommandRepository(db *gorm.DB) *sliderCommandRepository {
	return &sliderCommandRepository{db: db}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *sliderCommandRepository) Create(ctx context.Context, request *requests.CreateSliderRequest) (*models.Slider, error) {
	slider := &models.Slider{
		Name:      request.Nama,
		Image:     request.FilePath,
		CreatedAt: timePtr(time.Now()),
		UpdatedAt: timePtr(time.Now()),
	}

	if err := r.db.WithContext(ctx).Create(slider).Error; err != nil {
		return nil, slider_errors.ErrCreateSlider
	}
	return slider, nil
}

func (r *sliderCommandRepository) Update(ctx context.Context, request *requests.UpdateSliderRequest) (*models.Slider, error) {
	var slider models.Slider
	if err := r.db.WithContext(ctx).First(&slider, *request.ID).Error; err != nil {
		return nil, slider_errors.ErrSliderNotFound
	}

	slider.Name = request.Nama
	slider.Image = request.FilePath
	slider.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&slider).Error; err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}
	return &slider, nil
}

func (r *sliderCommandRepository) Trash(ctx context.Context, sliderID int) (*models.Slider, error) {
	var slider models.Slider
	if err := r.db.WithContext(ctx).First(&slider, sliderID).Error; err != nil {
		return nil, slider_errors.ErrSliderNotFound
	}
	now := timePtr(time.Now())
	if err := r.db.WithContext(ctx).Model(&slider).Update("deleted_at", now).Error; err != nil {
		return nil, slider_errors.ErrTrashSlider
	}
	slider.DeletedAt = now
	return &slider, nil
}

func (r *sliderCommandRepository) Restore(ctx context.Context, sliderID int) (*models.Slider, error) {
	var slider models.Slider
	if err := r.db.WithContext(ctx).Unscoped().Where("slider_id = ?", sliderID).First(&slider).Error; err != nil {
		return nil, slider_errors.ErrSliderNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&slider).Update("deleted_at", nil).Error; err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}
	slider.DeletedAt = nil
	return &slider, nil
}

func (r *sliderCommandRepository) DeletePermanent(ctx context.Context, sliderID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("slider_id = ?", sliderID).Delete(&models.Slider{})
	if result.Error != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}
	if result.RowsAffected == 0 {
		return false, slider_errors.ErrSliderNotFound
	}
	return true, nil
}

func (r *sliderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Model(&models.Slider{}).Update("deleted_at", nil)
	if result.Error != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Slider{})
	if result.Error != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
