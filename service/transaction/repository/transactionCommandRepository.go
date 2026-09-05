package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/transaction_errors"
	"gorm.io/gorm"
)

type transactionCommandRepository struct {
	db *gorm.DB
}

func NewTransactionCommandRepository(db *gorm.DB) *transactionCommandRepository {
	return &transactionCommandRepository{db: db}
}

func toTransactionResult(t *models.Transaction) *TransactionResult {
	if t == nil {
		return nil
	}
	return &TransactionResult{
		TransactionID: t.TransactionID,
		OrderID:       t.OrderID,
		MerchantID:    t.MerchantID,
		PaymentMethod: t.PaymentMethod,
		Amount:        t.Amount,
		PaymentStatus: t.PaymentStatus,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		DeletedAt:     t.DeletedAt,
	}
}

func (r *transactionCommandRepository) Create(ctx context.Context, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	now := time.Now()
	txn := &models.Transaction{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if err := r.db.WithContext(ctx).Create(txn).Error; err != nil {
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}
	return toTransactionResult(txn), nil
}

func (r *transactionCommandRepository) CreateInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateTransactionRequest) (*TransactionResult, error) {
	now := time.Now()
	txn := &models.Transaction{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if err := tx.WithContext(ctx).Create(txn).Error; err != nil {
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}
	return toTransactionResult(txn), nil
}

func (r *transactionCommandRepository) Update(ctx context.Context, request *requests.UpdateTransactionRequest) (*TransactionResult, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Where("transaction_id = ? AND deleted_at IS NULL", *request.TransactionID).First(&txn).Error; err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}

	txn.MerchantID = int32(request.MerchantID)
	txn.PaymentMethod = request.PaymentMethod
	txn.Amount = int32(request.Amount)
	txn.OrderID = int32(request.OrderID)
	txn.PaymentStatus = *request.PaymentStatus
	now := time.Now()
	txn.UpdatedAt = &now

	if err := r.db.WithContext(ctx).Save(&txn).Error; err != nil {
		return nil, transaction_errors.ErrUpdateTransaction.WithInternal(err)
	}

	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Trash(ctx context.Context, transaction_id int) (*TransactionResult, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Where("transaction_id = ? AND deleted_at IS NULL", transaction_id).First(&txn).Error; err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}

	now := time.Now()
	txn.DeletedAt = &now
	if err := r.db.WithContext(ctx).Save(&txn).Error; err != nil {
		return nil, transaction_errors.ErrTrashTransaction.WithInternal(err)
	}

	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) Restore(ctx context.Context, transaction_id int) (*TransactionResult, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Unscoped().Where("transaction_id = ? AND deleted_at IS NOT NULL", transaction_id).First(&txn).Error; err != nil {
		return nil, transaction_errors.ErrTransactionNotFound
	}

	txn.DeletedAt = nil
	if err := r.db.WithContext(ctx).Unscoped().Save(&txn).Error; err != nil {
		return nil, transaction_errors.ErrRestoreTransaction.WithInternal(err)
	}

	return toTransactionResult(&txn), nil
}

func (r *transactionCommandRepository) DeletePermanent(ctx context.Context, transaction_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("transaction_id = ?", transaction_id).Delete(&models.Transaction{})
	if result.Error != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, transaction_errors.ErrTransactionNotFound
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("order_id = ?", order_id).Delete(&models.Transaction{})
	if result.Error != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(result.Error)
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Transaction{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, transaction_errors.ErrRestoreAllTransactions.WithInternal(result.Error)
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Transaction{})
	if result.Error != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent.WithInternal(result.Error)
	}
	return true, nil
}
