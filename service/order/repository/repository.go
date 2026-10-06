package repository

import (
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc-pb/merchant"
	pborder_item "github.com/MamangRust/microservice-ecommerce-grpc-pb/order_item"
	pbproduct "github.com/MamangRust/microservice-ecommerce-grpc-pb/product"
	pbshipping_address "github.com/MamangRust/microservice-ecommerce-grpc-pb/shipping_address"
	pbtransaction "github.com/MamangRust/microservice-ecommerce-grpc-pb/transaction"
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc-pb/user"

	"github.com/MamangRust/microservice-ecommerce-pkg/adapter"
	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/order_item"
	productadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/product"
	shippingaddressadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/shipping_address"
	transactionadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/transaction"
	useradapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/user"
	"gorm.io/gorm"
)

// GuardOptions carries the resilience guard options for each outbound
// dependency.
type GuardOptions struct {
	User        []adapter.GuardOption
	Product     []adapter.GuardOption
	Merchant    []adapter.GuardOption
	OrderItem   []adapter.GuardOption
	Shipping    []adapter.GuardOption
	Transaction []adapter.GuardOption
}

type Repositories struct {
	OrderQuery       OrderQueryRepository
	OrderCommand     OrderCommandRepository
	UserQuery        UserQueryRepository
	ProductQuery     ProductQueryRepository
	ProductCommand   ProductCommandRepository
	OrderItemQuery   OrderItemQueryRepository
	OrderItemCommand OrderItemCommandRepository
	MerchantQuery    MerchantQueryRepository
	ShippingAddress  ShippingAddressCommandRepository
	TransactionCommand TransactionCommandRepository
	ShippingQuery    pbshipping_address.ShippingQueryServiceClient
	StockReservation StockReservationRepository
	Outbox           OutboxRepository
}

type Deps struct {
	GormDB *gorm.DB

	UserQueryClient      pbuser.UserQueryServiceClient
	ProductQueryClient   pbproduct.ProductQueryServiceClient
	ProductCommandClient pbproduct.ProductCommandServiceClient
	MerchantQueryClient  pbmerchant.MerchantQueryServiceClient
	OrderItemQueryClient pborder_item.OrderItemQueryServiceClient
	OrderItemCommandClient pborder_item.OrderItemCommandServiceClient
	ShippingCommandClient  pbshipping_address.ShippingCommandServiceClient
	ShippingQueryClient    pbshipping_address.ShippingQueryServiceClient
	TransactionCommandClient pbtransaction.TransactionCommandServiceClient

	Guards GuardOptions
}

func NewRepositories(deps *Deps) *Repositories {
	g := deps.Guards

	productAdapter := productadapter.New(deps.ProductQueryClient, deps.ProductCommandClient, g.Product...)
	orderItemAdapter := orderitemadapter.New(deps.OrderItemQueryClient, deps.OrderItemCommandClient, g.OrderItem...)
	shippingAdapter := shippingaddressadapter.New(deps.ShippingQueryClient, deps.ShippingCommandClient, g.Shipping...)

	return &Repositories{
		OrderQuery:         NewOrderQueryRepository(deps.GormDB),
		OrderCommand:       NewOrderCommandRepository(deps.GormDB),
		UserQuery:          useradapter.New(deps.UserQueryClient, nil, g.User...),
		ProductQuery:       productAdapter,
		ProductCommand:     productAdapter,
		OrderItemQuery:     orderItemAdapter,
		OrderItemCommand:   orderItemAdapter,
		MerchantQuery:      merchantadapter.New(deps.MerchantQueryClient, g.Merchant...),
		ShippingAddress:    shippingAdapter,
		TransactionCommand: transactionadapter.NewCommandAdapter(deps.TransactionCommandClient, g.Transaction...),
		ShippingQuery:      deps.ShippingQueryClient,
		StockReservation:   NewStockReservationRepository(deps.GormDB),
		Outbox:             NewOutboxRepository(deps.GormDB),
	}
}
