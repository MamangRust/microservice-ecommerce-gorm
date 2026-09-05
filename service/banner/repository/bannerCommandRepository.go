package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/banner_errors"
	"gorm.io/gorm"
)

type bannerCommandRepository struct {
	db *gorm.DB
}

func NewBannerCommandRepository(db *gorm.DB) *bannerCommandRepository {
	return &bannerCommandRepository{db: db}
}

func timePtr(t time.Time) *time.Time { return &t }

func parseDatePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	return &t, nil
}

func parseTimeStr(s string) (*string, error) {
	if s == "" {
		return nil, nil
	}
	_, err := time.Parse("15:04:05", s)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}
	return &s, nil
}

func (r *bannerCommandRepository) Create(ctx context.Context, request *requests.CreateBannerRequest) (*models.Banner, error) {
	startDate, err := parseDatePtr(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}
	endDate, err := parseDatePtr(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}
	startTime, err := parseTimeStr(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}
	endTime, err := parseTimeStr(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	banner := &models.Banner{
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  &request.IsActive,
		CreatedAt: timePtr(time.Now()),
		UpdatedAt: timePtr(time.Now()),
	}

	if err := r.db.WithContext(ctx).Create(banner).Error; err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}
	return banner, nil
}

func (r *bannerCommandRepository) Update(ctx context.Context, request *requests.UpdateBannerRequest) (*models.Banner, error) {
	startDate, err := parseDatePtr(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate.WithInternal(err)
	}
	endDate, err := parseDatePtr(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate.WithInternal(err)
	}
	startTime, err := parseTimeStr(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime.WithInternal(err)
	}
	endTime, err := parseTimeStr(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime.WithInternal(err)
	}

	var banner models.Banner
	if err := r.db.WithContext(ctx).First(&banner, *request.BannerID).Error; err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}

	banner.Name = request.Name
	banner.StartDate = startDate
	banner.EndDate = endDate
	banner.StartTime = startTime
	banner.EndTime = endTime
	banner.IsActive = &request.IsActive
	banner.UpdatedAt = timePtr(time.Now())

	if err := r.db.WithContext(ctx).Save(&banner).Error; err != nil {
		return nil, banner_errors.ErrUpdateBanner.WithInternal(err)
	}
	return &banner, nil
}

func (r *bannerCommandRepository) Trash(ctx context.Context, bannerID int) (*models.Banner, error) {
	var banner models.Banner
	if err := r.db.WithContext(ctx).First(&banner, bannerID).Error; err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}
	now := timePtr(time.Now())
	if err := r.db.WithContext(ctx).Model(&banner).Update("deleted_at", now).Error; err != nil {
		return nil, banner_errors.ErrTrashedBanner.WithInternal(err)
	}
	banner.DeletedAt = now
	return &banner, nil
}

func (r *bannerCommandRepository) Restore(ctx context.Context, bannerID int) (*models.Banner, error) {
	var banner models.Banner
	if err := r.db.WithContext(ctx).Unscoped().Where("banner_id = ?", bannerID).First(&banner).Error; err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&banner).Update("deleted_at", nil).Error; err != nil {
		return nil, banner_errors.ErrRestoreBanner.WithInternal(err)
	}
	banner.DeletedAt = nil
	return &banner, nil
}

func (r *bannerCommandRepository) DeletePermanent(ctx context.Context, bannerID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("banner_id = ?", bannerID).Delete(&models.Banner{})
	if result.Error != nil {
		return false, banner_errors.ErrDeleteBannerPermanent.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, banner_errors.ErrBannerNotFound
	}
	return true, nil
}

func (r *bannerCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Model(&models.Banner{}).Update("deleted_at", nil)
	if result.Error != nil {
		return false, banner_errors.ErrRestoreAllBanners.WithInternal(result.Error)
	}
	return true, nil
}

func (r *bannerCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Banner{})
	if result.Error != nil {
		return false, banner_errors.ErrDeleteAllBanners.WithInternal(result.Error)
	}
	return true, nil
}
