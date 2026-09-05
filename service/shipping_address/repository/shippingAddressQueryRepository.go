package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	shippingaddress_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/shipping_address_errors"
	"gorm.io/gorm"
)

type shippingAddressQueryRepository struct {
	db *gorm.DB
}

func NewShippingAddressQueryRepository(db *gorm.DB) *shippingAddressQueryRepository {
	return &shippingAddressQueryRepository{db: db}
}

func (r *shippingAddressQueryRepository) FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ShippingAddressResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT sa.shipping_address_id, sa.order_id, sa.alamat, sa.provinsi, sa.kota, sa.negara,
			sa.courier, sa.shipping_method, sa.shipping_cost,
			COALESCE(TO_CHAR(sa.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(sa.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(sa.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM shipping_addresses sa
		WHERE sa.deleted_at IS NULL
			AND (? = '' OR sa.alamat ILIKE ? OR sa.provinsi ILIKE ?)
		ORDER BY sa.shipping_address_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, shippingaddress_errors.ErrFindAllShippingAddress
	}
	return results, nil
}

func (r *shippingAddressQueryRepository) FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ShippingAddressResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT sa.shipping_address_id, sa.order_id, sa.alamat, sa.provinsi, sa.kota, sa.negara,
			sa.courier, sa.shipping_method, sa.shipping_cost,
			COALESCE(TO_CHAR(sa.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(sa.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(sa.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM shipping_addresses sa
		WHERE sa.deleted_at IS NULL
			AND (? = '' OR sa.alamat ILIKE ? OR sa.provinsi ILIKE ?)
		ORDER BY sa.shipping_address_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, shippingaddress_errors.ErrFindActiveShippingAddress
	}
	return results, nil
}

func (r *shippingAddressQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*ShippingAddressResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*ShippingAddressResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT sa.shipping_address_id, sa.order_id, sa.alamat, sa.provinsi, sa.kota, sa.negara,
			sa.courier, sa.shipping_method, sa.shipping_cost,
			COALESCE(TO_CHAR(sa.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(sa.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(sa.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at,
			COUNT(*) OVER() as total_count
		FROM shipping_addresses sa
		WHERE sa.deleted_at IS NOT NULL
			AND (? = '' OR sa.alamat ILIKE ? OR sa.provinsi ILIKE ?)
		ORDER BY sa.shipping_address_id DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, shippingaddress_errors.ErrFindTrashedShippingAddress
	}
	return results, nil
}

func (r *shippingAddressQueryRepository) FindByID(ctx context.Context, shippingID int) (*ShippingAddressResult, error) {
	var result ShippingAddressResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT sa.shipping_address_id, sa.order_id, sa.alamat, sa.provinsi, sa.kota, sa.negara,
			sa.courier, sa.shipping_method, sa.shipping_cost,
			COALESCE(TO_CHAR(sa.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(sa.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(sa.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at
		FROM shipping_addresses sa
		WHERE sa.shipping_address_id = ? AND sa.deleted_at IS NULL
	`, shippingID).Scan(&result).Error
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}
	if result.ShippingAddressID == 0 {
		return nil, shippingaddress_errors.ErrFindShippingAddressByID
	}
	return &result, nil
}

func (r *shippingAddressQueryRepository) FindByOrder(ctx context.Context, orderID int) (*ShippingAddressResult, error) {
	var result ShippingAddressResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT sa.shipping_address_id, sa.order_id, sa.alamat, sa.provinsi, sa.kota, sa.negara,
			sa.courier, sa.shipping_method, sa.shipping_cost,
			COALESCE(TO_CHAR(sa.created_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as created_at,
			COALESCE(TO_CHAR(sa.updated_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as updated_at,
			COALESCE(TO_CHAR(sa.deleted_at, 'YYYY-MM-DD HH24:MI:SS.MS'), '') as deleted_at
		FROM shipping_addresses sa
		WHERE sa.order_id = ?
	`, orderID).Scan(&result).Error
	if err != nil {
		return nil, shippingaddress_errors.ErrFindShippingAddressByOrder
	}
	return &result, nil
}
