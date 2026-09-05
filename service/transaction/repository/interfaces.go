package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-transaction/dto"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	"gorm.io/gorm"
)

type TransactionResult struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	PaymentStatus string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	TotalCount    int64
}

type UserQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*dto.GetUserByIDRow, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*dto.GetMerchantByIDRow, error)
}

type OrderItemRepository interface {
	FindOrderItemByOrder(
		ctx context.Context,
		order_id int,
	) ([]*dto.GetOrderItemsByOrderRow, error)
}

type OrderQueryRepository interface {
	FindByID(
		ctx context.Context,
		order_id int,
	) (*dto.GetOrderByIDRow, error)
}

type ShippingAddressQueryRepository interface {
	FindByID(
		ctx context.Context,
		shipping_id int,
	) (*dto.GetShippingAddressByOrderIDRow, error)
}

type TransactionQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllTransaction,
	) ([]*TransactionResult, error)

	FindByMerchant(
		ctx context.Context,
		req *requests.FindAllTransactionByMerchant,
	) ([]*TransactionResult, error)

	FindByID(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	FindByOrderID(
		ctx context.Context,
		order_id int,
	) (*TransactionResult, error)
}

type TransactionCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateTransactionRequest,
	) (*TransactionResult, error)

	CreateInTx(
		ctx context.Context,
		tx *gorm.DB,
		request *requests.CreateTransactionRequest,
	) (*TransactionResult, error)

	Update(
		ctx context.Context,
		request *requests.UpdateTransactionRequest,
	) (*TransactionResult, error)

	Trash(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	Restore(
		ctx context.Context,
		transaction_id int,
	) (*TransactionResult, error)

	DeletePermanent(
		ctx context.Context,
		transaction_id int,
	) (bool, error)

	DeleteByOrderIDPermanent(
		ctx context.Context,
		order_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*OutboxEventResult, error)
	CreateInTx(ctx context.Context, tx *gorm.DB, topic, key string, payload []byte) (*OutboxEventResult, error)
	GetPending(ctx context.Context, limit int) ([]*OutboxEventResult, error)
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*OutboxEventResult, error)
	MarkDelivered(ctx context.Context, outboxID int64) (*OutboxEventResult, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*OutboxEventResult, error)
	MarkDead(ctx context.Context, outboxID int64) (*OutboxEventResult, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type OutboxEventResult struct {
	OutboxID      int64
	Topic         string
	EventKey      string
	Payload       []byte
	Status        string
	Attempts      int32
	NextAttemptAt time.Time
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
