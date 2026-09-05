package repository

import (
	"context"

	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/transaction_errors"
	"gorm.io/gorm"
)

type transactionQueryRepository struct {
	db *gorm.DB
}

func NewTransactionQueryRepository(db *gorm.DB) *transactionQueryRepository {
	return &transactionQueryRepository{db: db}
}

func (r *transactionQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*TransactionResult
	err := r.db.WithContext(ctx).Model(&struct {
		TransactionID int32  `gorm:"column:transaction_id"`
		OrderID       int32  `gorm:"column:order_id"`
		MerchantID    int32  `gorm:"column:merchant_id"`
		PaymentMethod string `gorm:"column:payment_method"`
		Amount        int32  `gorm:"column:amount"`
		PaymentStatus string `gorm:"column:payment_status"`
		CreatedAt     *interface{}
		UpdatedAt     *interface{}
		DeletedAt     *interface{}
	}{}).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM transactions
		WHERE deleted_at IS NULL
			AND (? = '' OR payment_method ILIKE ? OR payment_status ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindAllTransactions.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindActive(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*TransactionResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM transactions
		WHERE deleted_at IS NULL
			AND (? = '' OR payment_method ILIKE ? OR payment_status ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindByActive.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*TransactionResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM transactions
		WHERE deleted_at IS NOT NULL
			AND (? = '' OR payment_method ILIKE ? OR payment_status ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindByTrashed.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*TransactionResult, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*TransactionResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM transactions
		WHERE deleted_at IS NULL
			AND merchant_id = ?
			AND (? = '' OR payment_method ILIKE ? OR payment_status ILIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, req.MerchantID, req.Search, "%" + req.Search + "%", "%" + req.Search + "%", req.PageSize, offset).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionQueryRepository) FindByID(ctx context.Context, transaction_id int) (*TransactionResult, error) {
	var result TransactionResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at, 0 AS total_count
		FROM transactions
		WHERE transaction_id = ? AND deleted_at IS NULL
	`, transaction_id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrFindById.WithInternal(err)
	}
	if result.TransactionID == 0 {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return &result, nil
}

func (r *transactionQueryRepository) FindByOrderID(ctx context.Context, order_id int) (*TransactionResult, error) {
	var result TransactionResult
	err := r.db.WithContext(ctx).Raw(`
		SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status,
			created_at, updated_at, deleted_at, 0 AS total_count
		FROM transactions
		WHERE order_id = ? AND deleted_at IS NULL
	`, order_id).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrFindByOrderId.WithInternal(err)
	}
	if result.TransactionID == 0 {
		return nil, transaction_errors.ErrTransactionNotFound
	}
	return &result, nil
}
