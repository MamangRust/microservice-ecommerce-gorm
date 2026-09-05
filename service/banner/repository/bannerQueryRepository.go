package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/banner_errors"
	"gorm.io/gorm"
)

type bannerQueryRepository struct {
	db *gorm.DB
}

func NewBannerQueryRepository(db *gorm.DB) *bannerQueryRepository {
	return &bannerQueryRepository{db: db}
}

func (r *bannerQueryRepository) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*BannerResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT b.banner_id, b.name,
			COALESCE(TO_CHAR(b.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(TO_CHAR(b.end_date, 'YYYY-MM-DD'), '') as end_date,
			COALESCE(TO_CHAR(b.start_time, 'HH24:MI:SS'), '') as start_time,
			COALESCE(TO_CHAR(b.end_time, 'HH24:MI:SS'), '') as end_time,
			b.is_active,
			COALESCE(TO_CHAR(b.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(b.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(b.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM banners b
		WHERE b.deleted_at IS NULL
			AND (? = '' OR b.name ILIKE ?)
		ORDER BY b.banner_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	return results, nil
}

func (r *bannerQueryRepository) FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*BannerResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT b.banner_id, b.name,
			COALESCE(TO_CHAR(b.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(TO_CHAR(b.end_date, 'YYYY-MM-DD'), '') as end_date,
			COALESCE(TO_CHAR(b.start_time, 'HH24:MI:SS'), '') as start_time,
			COALESCE(TO_CHAR(b.end_time, 'HH24:MI:SS'), '') as end_time,
			b.is_active,
			COALESCE(TO_CHAR(b.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(b.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(b.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM banners b
		WHERE b.deleted_at IS NULL
			AND b.is_active = true
			AND (? = '' OR b.name ILIKE ?)
		ORDER BY b.banner_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, banner_errors.ErrFindActiveBanners.WithInternal(err)
	}
	return results, nil
}

func (r *bannerQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*BannerResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*BannerResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT b.banner_id, b.name,
			COALESCE(TO_CHAR(b.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(TO_CHAR(b.end_date, 'YYYY-MM-DD'), '') as end_date,
			COALESCE(TO_CHAR(b.start_time, 'HH24:MI:SS'), '') as start_time,
			COALESCE(TO_CHAR(b.end_time, 'HH24:MI:SS'), '') as end_time,
			b.is_active,
			COALESCE(TO_CHAR(b.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(b.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(b.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM banners b
		WHERE b.deleted_at IS NOT NULL
			AND (? = '' OR b.name ILIKE ?)
		ORDER BY b.banner_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, banner_errors.ErrFindTrashedBanners.WithInternal(err)
	}
	return results, nil
}

func (r *bannerQueryRepository) FindByID(ctx context.Context, bannerID int) (*BannerResult, error) {
	var result BannerResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT b.banner_id, b.name,
			COALESCE(TO_CHAR(b.start_date, 'YYYY-MM-DD'), '') as start_date,
			COALESCE(TO_CHAR(b.end_date, 'YYYY-MM-DD'), '') as end_date,
			COALESCE(TO_CHAR(b.start_time, 'HH24:MI:SS'), '') as start_time,
			COALESCE(TO_CHAR(b.end_time, 'HH24:MI:SS'), '') as end_time,
			b.is_active,
			COALESCE(TO_CHAR(b.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(b.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(b.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at
		FROM banners b
		WHERE b.banner_id = ? AND b.deleted_at IS NULL
	`, bannerID).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, banner_errors.ErrBannerNotFound
		}
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}
	if result.BannerID == 0 {
		return nil, banner_errors.ErrBannerNotFound
	}
	return &result, nil
}
