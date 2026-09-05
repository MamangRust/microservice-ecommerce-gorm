package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/cart_errors"
	"gorm.io/gorm"
)

type cartCommandRepository struct {
	db *gorm.DB
}

func NewCartCommandRepository(db *gorm.DB) CartCommandRepository {
	return &cartCommandRepository{db: db}
}

func (r *cartCommandRepository) CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*CartCreateResult, error) {
	now := time.Now()

	type row struct {
		CartID    int32
		UserID    int32
		ProductID int32
		Name      string
		Price     int32
		Image     string
		Quantity  int32
		Weight    int32
		CreatedAt string
		UpdatedAt string
	}

	var res row
	err := r.db.WithContext(ctx).Raw(`
		INSERT INTO carts (user_id, product_id, name, price, image, quantity, weight, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING cart_id, user_id, product_id, name, price, image, quantity, weight,
		          TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS'),
		          TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')`,
		req.UserID, req.ProductID, req.Name, req.Price, req.ImageProduct, req.Quantity, req.Weight, now, now,
	).Scan(&res).Error
	if err != nil {
		return nil, cart_errors.ErrCreateCart.WithInternal(err)
	}

	return &CartCreateResult{
		CartID: res.CartID, UserID: res.UserID, ProductID: res.ProductID,
		Name: res.Name, Price: res.Price, Image: res.Image,
		Quantity: res.Quantity, Weight: res.Weight,
		CreatedAt: res.CreatedAt, UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *cartCommandRepository) DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error) {
	result := r.db.WithContext(ctx).Exec(
		`DELETE FROM carts WHERE cart_id = ? AND user_id = ?`, req.CartID, req.UserID,
	)
	if result.Error != nil {
		return false, cart_errors.ErrDeleteCartPermanent.WithInternal(result.Error)
	}
	return true, nil
}

func (r *cartCommandRepository) DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error) {
	cartIDs := make([]int, len(req.CartIds))
	for i, id := range req.CartIds {
		cartIDs[i] = id
	}

	result := r.db.WithContext(ctx).Exec(
		`DELETE FROM carts WHERE cart_id IN ? AND user_id = ?`, cartIDs, req.UserID,
	)
	if result.Error != nil {
		return false, cart_errors.ErrDeleteAllCarts.WithInternal(result.Error)
	}
	return true, nil
}
