package repository

import (
	"context"
	"time"

	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type OrderResult struct {
	OrderID    int32
	UserID     int32
	MerchantID int32
	TotalPrice int32
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type StockReservationResult struct {
	ReservationID int32
	OrderID       int32
	ProductID     int32
	Quantity      int32
	Status        string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// UserQueryRepository is the shared user query contract, provided by the user
// gRPC adapter.
type UserQueryRepository = useradapter.QueryRepository

// ProductQueryRepository is the shared product query contract, provided by the
// product gRPC adapter.
type ProductQueryRepository = productadapter.QueryRepository

// MerchantQueryRepository is the shared merchant query contract, provided by the
// merchant gRPC adapter so this service never holds a raw gRPC client.
type MerchantQueryRepository = merchantadapter.QueryRepository

// ProductCommandRepository is the shared product command contract, provided by
// the product gRPC adapter.
type ProductCommandRepository = productadapter.CommandRepository

// ShippingAddressCommandRepository is the shared shipping-address command
// contract, provided by the shipping-address gRPC adapter.
type ShippingAddressCommandRepository = shippingaddressadapter.CommandRepository

// TransactionCommandRepository is the shared transaction command contract,
// provided by the transaction gRPC adapter.
type TransactionCommandRepository = transactionadapter.CommandRepository

// OrderItemQueryRepository is the shared order-item query contract plus the
// price calculation the order flow needs.
type OrderItemQueryRepository interface {
	orderitemadapter.QueryRepository
	CalculateTotalPrice(ctx context.Context, order_id int) (*int32, error)
}

// OrderItemCommandRepository is the shared order-item command contract, provided
// by the order-item gRPC adapter.
type OrderItemCommandRepository = orderitemadapter.CommandRepository

type OrderCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error)
	Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error)
	Trash(ctx context.Context, order_id int) (*models.Order, error)
	Restore(ctx context.Context, order_id int) (*models.Order, error)
	DeletePermanent(ctx context.Context, order_id int) (bool, error)
	DeletePermanentWithChildren(ctx context.Context, order_id int) (bool, error)
	FindTrashedByID(ctx context.Context, order_id int) (*models.Order, error)
	FindTrashed(ctx context.Context) ([]*models.Order, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type OrderQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*OrderResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*OrderResult, error)
	FindByID(ctx context.Context, order_id int) (*models.Order, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*models.OutboxEvent, error)
	GetPending(ctx context.Context, limit int) ([]*models.OutboxEvent, error)
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*models.OutboxEvent, error)
	MarkDelivered(ctx context.Context, outboxID int64) (*models.OutboxEvent, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*models.OutboxEvent, error)
	MarkDead(ctx context.Context, outboxID int64) (*models.OutboxEvent, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type StockReservationRepository interface {
	GetByOrder(ctx context.Context, orderID int) ([]*models.OrderStockReservation, error)
	Upsert(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error)
	UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*models.OrderStockReservation, error)
	Release(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error)
	Reserve(ctx context.Context, orderID, productID int) (*models.OrderStockReservation, error)
	GetReservedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	GetReleasedForTrashedOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	DeleteByOrder(ctx context.Context, orderID int) error
	DeleteByOrderProduct(ctx context.Context, orderID, productID int) error
	DeleteAllForTrashedOrders(ctx context.Context) error
	GetReleasedForActiveOrders(ctx context.Context) ([]*models.OrderStockReservation, error)
	DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error)
}
