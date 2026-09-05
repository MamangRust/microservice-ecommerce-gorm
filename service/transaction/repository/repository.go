package repository

import (
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc/pb/order_item"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc/pb/shipping_address"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	pborder "github.com/MamangRust/microservice-ecommerce-grpc/pb/order"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	TransactionCommand TransactionCommandRepository
	TransactionQuery   TransactionQueryRepository
	OrderItem          OrderItemRepository
	OrderQuery         OrderQueryRepository
	MerchantQuery      MerchantQueryRepository
	ShippingAddress    ShippingAddressQueryRepository
	UserQuery          UserQueryRepository
	Outbox             OutboxRepository
}

type Deps struct {
	DB             *gorm.DB
	UserQuery      pbuser.UserQueryServiceClient
	MerchantQuery  pbmerchant.MerchantQueryServiceClient
	OrderQuery     pborder.OrderQueryServiceClient
	OrderItemQuery pborder_item.OrderItemQueryServiceClient
	ShippingQuery  pbshipping_address.ShippingQueryServiceClient
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.DB),
		TransactionQuery:   NewTransactionQueryRepository(deps.DB),
		OrderItem:          NewOrderItemRepository(deps.OrderItemQuery),
		OrderQuery:         NewOrderQueryRepository(deps.OrderQuery),
		MerchantQuery:      NewMerchantQueryRepository(deps.MerchantQuery),
		ShippingAddress:    NewShippingAddressQueryRepository(deps.ShippingQuery),
		UserQuery:          NewUserQueryRepository(deps.UserQuery),
		Outbox:             NewOutboxRepository(deps.DB),
	}
}
