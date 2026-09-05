package repository

import (
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc/pb/product"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc/pb/transaction"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantQuery        MerchantQueryRepository
	ProductQuery         ProductQueryRepository
	ProductCommand       ProductCommandRepository
	OrderItemQuery       OrderItemQueryRepository
	OrderItemCommand     OrderItemCommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	UserQuery            UserQueryRepository
	ShippingAddress      ShippingAddressCommandRepository
	TransactionCommand   TransactionCommandRepository
	ShippingQuery        pbshipping_address.ShippingQueryServiceClient
	StockReservation     StockReservationRepository
	Outbox               OutboxRepository
}

type Deps struct {
	DB               *gorm.DB
	MerchantQuery    pbmerchant.MerchantQueryServiceClient
	ProductQuery     pbproduct.ProductQueryServiceClient
	ProductCommand   pbproduct.ProductCommandServiceClient
	OrderItemQuery   pborder_item.OrderItemQueryServiceClient
	OrderItemCommand pborder_item.OrderItemCommandServiceClient
	UserQuery        pbuser.UserQueryServiceClient
	ShippingCommand  pbshipping_address.ShippingCommandServiceClient
	TransactionCommand pbtransaction.TransactionCommandServiceClient
	ShippingQuery      pbshipping_address.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		MerchantQuery:    NewMerchantQueryRepository(deps.MerchantQuery),
		ProductQuery:     NewProductQueryRepository(deps.ProductQuery),
		ProductCommand:   NewProductCommandRepository(deps.ProductCommand),
		OrderItemQuery:   NewOrderItemQueryRepository(deps.OrderItemQuery, deps.OrderItemCommand),
		OrderItemCommand: NewOrderItemCommandRepository(deps.OrderItemCommand),
		OrderQuery:       NewOrderQueryRepository(deps.DB),
		OrderCommand:     NewOrderCommandRepository(deps.DB),
		UserQuery:        NewUserQueryRepository(deps.UserQuery),
		ShippingAddress:  NewShippingAddressCommandRepository(deps.ShippingCommand),
		TransactionCommand: NewTransactionCommandRepository(deps.TransactionCommand),
		ShippingQuery:      deps.ShippingQuery,
		StockReservation:   NewStockReservationRepository(deps.DB),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
