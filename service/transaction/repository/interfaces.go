package repository

import (
	"context"
	"time"

	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
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

// UserQueryRepository is the shared user query contract, provided by the user
// gRPC adapter.
type UserQueryRepository = useradapter.QueryRepository

// MerchantQueryRepository is the shared merchant query contract, provided by the
// merchant gRPC adapter so this service never holds a raw gRPC client.
type MerchantQueryRepository = merchantadapter.QueryRepository

// OrderItemRepository is the shared order-item query contract, provided by the
// order-item gRPC adapter.
type OrderItemRepository = orderitemadapter.QueryRepository

// OrderQueryRepository is the shared order query contract, provided by the
// order gRPC adapter.
type OrderQueryRepository = orderadapter.QueryRepository

// ShippingAddressQueryRepository is the shared shipping-address query contract,
// provided by the shipping-address gRPC adapter.
type ShippingAddressQueryRepository = shippingaddressadapter.QueryRepository

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
