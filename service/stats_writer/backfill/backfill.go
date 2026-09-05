// Package backfill implements the stats-writer `backfill` command: it reads
// historical OLTP rows from the per-service PostgreSQL databases (orders,
// order_items joined with products/categories, transactions) and materializes
// them into ClickHouse through the same batch repository used for live events.
//
// This is the bootstrap path for the stats pipeline — it lets the ClickHouse
// tables reflect pre-existing data without replaying every domain event.
package backfill

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/microservice-ecommerce-grpc-stats-writer/repository"
	"github.com/MamangRust/microservice-ecommerce-pkg/database"
	"github.com/MamangRust/microservice-ecommerce-pkg/logger"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// backfillEventID derives a deterministic UUID per entity so re-running the
// backfill replaces the same ReplacingMergeTree key (with a newer version)
// instead of appending duplicates.
func backfillEventID(kind string, id int32) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("backfill:%s:%d", kind, id))).String()
}

// Backfiller reads OLTP rows and pushes them into ClickHouse.
type Backfiller struct {
	log      logger.LoggerInterface
	repo     repository.Repository
	order    *gorm.DB
	item     *gorm.DB
	product  *gorm.DB
	category *gorm.DB
	tx       *gorm.DB
}

// New opens one GORM connection per service database that owns a stats source and
// returns a ready Backfiller. Call Close to release them.
func New(log logger.LoggerInterface, repo repository.Repository) (*Backfiller, error) {
	open := func(prefix string) (*gorm.DB, error) {
		conn, err := database.NewGormClientWithPrefix(log, prefix)
		if err != nil {
			return nil, fmt.Errorf("connect %s: %w", prefix, err)
		}
		return conn, nil
	}

	order, err := open("DB_ORDER")
	if err != nil {
		return nil, err
	}
	item, err := open("DB_ORDER_ITEM")
	if err != nil {
		return nil, err
	}
	product, err := open("DB_PRODUCT")
	if err != nil {
		return nil, err
	}
	category, err := open("DB_CATEGORY")
	if err != nil {
		return nil, err
	}
	tx, err := open("DB_TRANSACTION")
	if err != nil {
		return nil, err
	}

	return &Backfiller{
		log:      log,
		repo:     repo,
		order:    order,
		item:     item,
		product:  product,
		category: category,
		tx:       tx,
	}, nil
}

func (b *Backfiller) Close() {
	for _, conn := range []*gorm.DB{b.order, b.item, b.product, b.category, b.tx} {
		if conn != nil {
			if sqlDB, err := conn.DB(); err == nil {
				sqlDB.Close()
			}
		}
	}
}

// Run streams all stats sources into ClickHouse. The event version is the
// backfill run timestamp so re-running supersedes previous rows.
func (b *Backfiller) Run(ctx context.Context) error {
	version := uint64(time.Now().Unix())
	counts := map[string]int{}

	if err := b.backfillOrders(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillOrderItems(ctx, version, counts); err != nil {
		return err
	}
	if err := b.backfillTransactions(ctx, version, counts); err != nil {
		return err
	}

	if err := b.repo.Flush(ctx); err != nil {
		return fmt.Errorf("flush backfill batches: %w", err)
	}

	b.log.Info("backfill complete",
		zap.Int("orders", counts["order"]),
		zap.Int("order_items", counts["order_item"]),
		zap.Int("transactions", counts["transaction"]),
	)
	return nil
}

type orderRow struct {
	OrderID    int32
	UserID     int32
	MerchantID int32
	TotalPrice int32
	CreatedAt  time.Time
}

type itemRow struct {
	OrderItemID int32
	OrderID     int32
	ProductID   int32
	Quantity    int32
	Price       int32
	CreatedAt   time.Time
}

type txRow struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	Status        string
	CreatedAt     time.Time
}

type orderMerchantRow struct {
	OrderID    int32
	MerchantID int32
}

type productCategoryRow struct {
	ProductID  int32
	CategoryID int32
}

type categoryNameRow struct {
	CategoryID int32
	Name       string
}

func (b *Backfiller) backfillOrders(ctx context.Context, version uint64, counts map[string]int) error {
	var rows []orderRow
	if err := b.order.WithContext(ctx).Raw(`SELECT order_id, user_id, merchant_id, total_price, created_at FROM orders WHERE deleted_at IS NULL`).Scan(&rows).Error; err != nil {
		return fmt.Errorf("query orders: %w", err)
	}

	for _, r := range rows {
		event := events.OrderEvent{
			OrderID:    r.OrderID,
			UserID:     r.UserID,
			MerchantID: r.MerchantID,
			TotalPrice: r.TotalPrice,
			Status:     "created",
			EventTime:  r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderEvent(ctx, backfillEventID("order", r.OrderID), version, event); err != nil {
			return fmt.Errorf("insert order %d: %w", r.OrderID, err)
		}
		counts["order"]++
	}
	return nil
}

func (b *Backfiller) backfillOrderItems(ctx context.Context, version uint64, counts map[string]int) error {
	// Load order -> merchant_id from the order service DB.
	var orderMerchantRows []orderMerchantRow
	if err := b.order.WithContext(ctx).Raw(`SELECT order_id, merchant_id FROM orders WHERE deleted_at IS NULL`).Scan(&orderMerchantRows).Error; err != nil {
		return fmt.Errorf("query order merchant map: %w", err)
	}
	orderMerchant := map[int32]int32{}
	for _, r := range orderMerchantRows {
		orderMerchant[r.OrderID] = r.MerchantID
	}

	// Load product -> category_id from the product service DB.
	var productCategoryRows []productCategoryRow
	if err := b.product.WithContext(ctx).Raw(`SELECT product_id, category_id FROM products WHERE deleted_at IS NULL`).Scan(&productCategoryRows).Error; err != nil {
		return fmt.Errorf("query product category map: %w", err)
	}
	productCategory := map[int32]int32{}
	for _, r := range productCategoryRows {
		productCategory[r.ProductID] = r.CategoryID
	}

	// Load category_id -> name from the category service DB.
	var categoryNameRows []categoryNameRow
	if err := b.category.WithContext(ctx).Raw(`SELECT category_id, name FROM categories WHERE deleted_at IS NULL`).Scan(&categoryNameRows).Error; err != nil {
		return fmt.Errorf("query category name map: %w", err)
	}
	categoryName := map[int32]string{}
	for _, r := range categoryNameRows {
		categoryName[r.CategoryID] = r.Name
	}

	var itemRows []itemRow
	if err := b.item.WithContext(ctx).Raw(`SELECT order_item_id, order_id, product_id, quantity, price, created_at FROM order_items WHERE deleted_at IS NULL`).Scan(&itemRows).Error; err != nil {
		return fmt.Errorf("query order items: %w", err)
	}

	for _, r := range itemRows {
		catID := productCategory[r.ProductID]
		event := events.OrderItemEvent{
			OrderItemID:  r.OrderItemID,
			OrderID:      r.OrderID,
			MerchantID:   orderMerchant[r.OrderID],
			ProductID:    r.ProductID,
			CategoryID:   catID,
			CategoryName: categoryName[catID],
			Quantity:     r.Quantity,
			Price:        r.Price,
			EventTime:    r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertOrderItemEvent(ctx, backfillEventID("order_item", r.OrderItemID), version, event); err != nil {
			return fmt.Errorf("insert order item %d: %w", r.OrderItemID, err)
		}
		counts["order_item"]++
	}
	return nil
}

func (b *Backfiller) backfillTransactions(ctx context.Context, version uint64, counts map[string]int) error {
	var rows []txRow
	if err := b.tx.WithContext(ctx).Raw(`SELECT transaction_id, order_id, merchant_id, payment_method, amount, payment_status, created_at FROM transactions WHERE deleted_at IS NULL`).Scan(&rows).Error; err != nil {
		return fmt.Errorf("query transactions: %w", err)
	}

	for _, r := range rows {
		event := events.TransactionEvent{
			TransactionID: r.TransactionID,
			OrderID:       r.OrderID,
			MerchantID:    r.MerchantID,
			PaymentMethod: r.PaymentMethod,
			Amount:        r.Amount,
			Status:        r.Status,
			EventTime:     r.CreatedAt.UTC().Format(time.RFC3339),
		}
		if err := b.repo.InsertTransactionEvent(ctx, backfillEventID("transaction", r.TransactionID), version, event); err != nil {
			return fmt.Errorf("insert transaction %d: %w", r.TransactionID, err)
		}
		counts["transaction"]++
	}
	return nil
}
