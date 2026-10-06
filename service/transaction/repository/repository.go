package repository

import (
	pborder_items "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	pborders "github.com/MamangRust/microservice-ecommerce-grpc-pb/order"
	pbmerchants "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pbshipping_addresses "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	pbusers "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"
	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
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

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User      []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
	Shipping  []adapter.GuardOption
}

type Deps struct {
	GormDB *gorm.DB

	UserQueryClient      pbusers.UserQueryServiceClient
	MerchantQueryClient  pbmerchants.MerchantQueryServiceClient
	OrderQueryClient     pborders.OrderQueryServiceClient
	OrderItemQueryClient pborder_items.OrderItemQueryServiceClient
	ShippingQueryClient  pbshipping_addresses.ShippingQueryServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(deps.GormDB),
		TransactionQuery:   NewTransactionQueryRepository(deps.GormDB),
		OrderItem:          orderitemadapter.New(deps.OrderItemQueryClient, nil, g.OrderItem...),
		OrderQuery:         orderadapter.New(deps.OrderQueryClient, g.Order...),
		MerchantQuery:      merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ShippingAddress:    shippingaddressadapter.New(deps.ShippingQueryClient, nil, g.Shipping...),
		UserQuery:          useradapter.New(deps.UserQueryClient, nil, g.User...),
		Outbox:             NewOutboxRepository(deps.GormDB),
	}
}
