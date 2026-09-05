package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/cart_errors"
	"gorm.io/gorm"
)

type cartQueryRepository struct {
	db *gorm.DB
}

func NewCartQueryRepository(db *gorm.DB) CartQueryRepository {
	return &cartQueryRepository{db: db}
}

func (r *cartQueryRepository) FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*CartResult, error) {

	type row struct {
		CartID     int32
		UserID     int32
		ProductID  int32
		Name       string
		Price      int32
		Image      string
		Quantity   int32
		Weight     int32
		CreatedAt  string
		UpdatedAt  string
		TotalCount int64
	}

	var results []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.cart_id, c.user_id, c.product_id, c.name, c.price, c.image,
		       c.quantity, c.weight,
		       TO_CHAR(c.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at,
		       TO_CHAR(c.updated_at, 'YYYY-MM-DD HH24:MI:SS') AS updated_at,
		       COUNT(*) OVER() AS total_count
		FROM carts c
		WHERE c.user_id = ? AND c.deleted_at IS NULL
		  AND (? = '' OR c.name ILIKE ?)
		ORDER BY c.created_at DESC
	`, req.UserID, req.Search, "%" + req.Search + "%").Scan(&results).Error
	if err != nil {
		return nil, cart_errors.ErrFindAllCarts.WithInternal(err)
	}

	var out []*CartResult
	for _, r := range results {
		cr := r
		out = append(out, &CartResult{
			CartID: cr.CartID, UserID: cr.UserID, ProductID: cr.ProductID,
			Name: cr.Name, Price: cr.Price, Image: cr.Image,
			Quantity: cr.Quantity, Weight: cr.Weight,
			CreatedAt: cr.CreatedAt, UpdatedAt: cr.UpdatedAt,
			TotalCount: cr.TotalCount,
		})
	}
	return out, nil
}
