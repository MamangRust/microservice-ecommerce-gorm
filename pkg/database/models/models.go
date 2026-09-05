package models

import (
	"time"
)

// ---------------------------------------------------------------------------
// Core
// ---------------------------------------------------------------------------

type User struct {
	UserID           int32      `gorm:"column:user_id;primaryKey" json:"user_id"`
	Firstname        string     `gorm:"column:firstname" json:"firstname"`
	Lastname         string     `gorm:"column:lastname" json:"lastname"`
	Email            string     `gorm:"column:email" json:"email"`
	Password         string     `gorm:"column:password" json:"password"`
	VerificationCode string     `gorm:"column:verification_code" json:"verification_code"`
	IsVerified       *bool      `gorm:"column:is_verified" json:"is_verified"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (User) TableName() string { return "users" }

type Role struct {
	RoleID    int32      `gorm:"column:role_id;primaryKey" json:"role_id"`
	RoleName  string     `gorm:"column:role_name" json:"role_name"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Role) TableName() string { return "roles" }

type UserRole struct {
	UserRoleID int32      `gorm:"column:user_role_id;primaryKey" json:"user_role_id"`
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	RoleID     int32      `gorm:"column:role_id" json:"role_id"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// ---------------------------------------------------------------------------
// Merchant
// ---------------------------------------------------------------------------

type Merchant struct {
	MerchantID   int32      `gorm:"column:merchant_id;primaryKey" json:"merchant_id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Address      *string    `gorm:"column:address" json:"address"`
	ContactEmail *string    `gorm:"column:contact_email" json:"contact_email"`
	ContactPhone *string    `gorm:"column:contact_phone" json:"contact_phone"`
	Status       string     `gorm:"column:status" json:"status"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Merchant) TableName() string { return "merchants" }

type MerchantDocument struct {
	DocumentID   int32      `gorm:"column:document_id;primaryKey" json:"document_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string     `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string     `gorm:"column:document_url" json:"document_url"`
	Status       string     `gorm:"column:status" json:"status"`
	Note         *string    `gorm:"column:note" json:"note"`
	UploadedAt   *time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantDocument) TableName() string { return "merchant_documents" }

type MerchantDetail struct {
	MerchantDetailID int32      `gorm:"column:merchant_detail_id;primaryKey" json:"merchant_detail_id"`
	MerchantID       int32      `gorm:"column:merchant_id" json:"merchant_id"`
	DisplayName      *string    `gorm:"column:display_name" json:"display_name"`
	CoverImageUrl    *string    `gorm:"column:cover_image_url" json:"cover_image_url"`
	LogoUrl          *string    `gorm:"column:logo_url" json:"logo_url"`
	ShortDescription *string    `gorm:"column:short_description" json:"short_description"`
	WebsiteUrl       *string    `gorm:"column:website_url" json:"website_url"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantDetail) TableName() string { return "merchant_details" }

type MerchantSocialMediaLink struct {
	MerchantSocialID int32      `gorm:"column:merchant_social_id;primaryKey" json:"merchant_social_id"`
	MerchantDetailID int32      `gorm:"column:merchant_detail_id" json:"merchant_detail_id"`
	Platform         string     `gorm:"column:platform" json:"platform"`
	Url              string     `gorm:"column:url" json:"url"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantSocialMediaLink) TableName() string { return "merchant_social_media_links" }

type MerchantBusinessInformation struct {
	MerchantBusinessInfoID int32      `gorm:"column:merchant_business_info_id;primaryKey" json:"merchant_business_info_id"`
	MerchantID             int32      `gorm:"column:merchant_id" json:"merchant_id"`
	BusinessType           *string    `gorm:"column:business_type" json:"business_type"`
	TaxID                  *string    `gorm:"column:tax_id" json:"tax_id"`
	EstablishedYear        *int32     `gorm:"column:established_year" json:"established_year"`
	NumberOfEmployees      *int32     `gorm:"column:number_of_employees" json:"number_of_employees"`
	WebsiteUrl             *string    `gorm:"column:website_url" json:"website_url"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt              *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantBusinessInformation) TableName() string { return "merchant_business_information" }

type MerchantPolicy struct {
	MerchantPolicyID int32      `gorm:"column:merchant_policy_id;primaryKey" json:"merchant_policy_id"`
	MerchantID       int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PolicyType       string     `gorm:"column:policy_type" json:"policy_type"`
	Title            string     `gorm:"column:title" json:"title"`
	Description      string     `gorm:"column:description" json:"description"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantPolicy) TableName() string { return "merchant_policies" }

type MerchantCertificationsAndAward struct {
	MerchantCertificationID int32      `gorm:"column:merchant_certification_id;primaryKey" json:"merchant_certification_id"`
	MerchantID              int32      `gorm:"column:merchant_id" json:"merchant_id"`
	Title                   string     `gorm:"column:title" json:"title"`
	Description             *string    `gorm:"column:description" json:"description"`
	IssuedBy                *string    `gorm:"column:issued_by" json:"issued_by"`
	IssueDate               *time.Time `gorm:"column:issue_date" json:"issue_date"`
	ExpiryDate              *time.Time `gorm:"column:expiry_date" json:"expiry_date"`
	CertificateUrl          *string    `gorm:"column:certificate_url" json:"certificate_url"`
	CreatedAt               *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt               *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantCertificationsAndAward) TableName() string { return "merchant_certifications_and_awards" }

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

type Product struct {
	ProductID    int32      `gorm:"column:product_id;primaryKey" json:"product_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Price        int32      `gorm:"column:price" json:"price"`
	CountInStock int32      `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string    `gorm:"column:brand" json:"brand"`
	Weight       *int32     `gorm:"column:weight" json:"weight"`
	Rating       *float64   `gorm:"column:rating" json:"rating"`
	SlugProduct  *string    `gorm:"column:slug_product" json:"slug_product"`
	ImageProduct *string    `gorm:"column:image_product" json:"image_product"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Product) TableName() string { return "products" }

type ProductStockAdjustment struct {
	OperationID string     `gorm:"column:operation_id;primaryKey" json:"operation_id"`
	ProductID   int32      `gorm:"column:product_id" json:"product_id"`
	Delta       int32      `gorm:"column:delta" json:"delta"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ProductStockAdjustment) TableName() string { return "product_stock_adjustments" }

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

type Order struct {
	OrderID    int32      `gorm:"column:order_id;primaryKey" json:"order_id"`
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	TotalPrice int32      `gorm:"column:total_price" json:"total_price"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Order) TableName() string { return "orders" }

type OrderItem struct {
	OrderItemID int32      `gorm:"column:order_item_id;primaryKey" json:"order_item_id"`
	OrderID     int32      `gorm:"column:order_id" json:"order_id"`
	ProductID   int32      `gorm:"column:product_id" json:"product_id"`
	Quantity    int32      `gorm:"column:quantity" json:"quantity"`
	Price       int32      `gorm:"column:price" json:"price"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OrderItem) TableName() string { return "order_items" }

type OrderStockReservation struct {
	ReservationID int32      `gorm:"column:reservation_id;primaryKey" json:"reservation_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	ProductID     int32      `gorm:"column:product_id" json:"product_id"`
	Quantity      int32      `gorm:"column:quantity" json:"quantity"`
	Status        string     `gorm:"column:status" json:"status"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (OrderStockReservation) TableName() string { return "order_stock_reservations" }

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

type RefreshToken struct {
	RefreshTokenID int32      `gorm:"column:refresh_token_id;primaryKey" json:"refresh_token_id"`
	UserID         int32      `gorm:"column:user_id" json:"user_id"`
	Token          string     `gorm:"column:token" json:"token"`
	Expiration     time.Time  `gorm:"column:expiration" json:"expiration"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

type ResetToken struct {
	ResetTokenID int32     `gorm:"column:id;primaryKey" json:"id"`
	UserID       int64     `gorm:"column:user_id" json:"user_id"`
	Token        string    `gorm:"column:token" json:"token"`
	ExpiryDate   time.Time `gorm:"column:expiry_date" json:"expiry_date"`
}

func (ResetToken) TableName() string { return "reset_tokens" }

// ---------------------------------------------------------------------------
// Content
// ---------------------------------------------------------------------------

type Banner struct {
	BannerID  int32      `gorm:"column:banner_id;primaryKey" json:"banner_id"`
	Name      string     `gorm:"column:name" json:"name"`
	StartDate *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate   *time.Time `gorm:"column:end_date" json:"end_date"`
	StartTime *string `gorm:"column:start_time" json:"start_time"`
	EndTime   *string `gorm:"column:end_time" json:"end_time"`
	IsActive  *bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Banner) TableName() string { return "banners" }

type Slider struct {
	SliderID  int32      `gorm:"column:slider_id;primaryKey" json:"slider_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Image     string     `gorm:"column:image" json:"image"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Slider) TableName() string { return "sliders" }

type Category struct {
	CategoryID    int32      `gorm:"column:category_id;primaryKey" json:"category_id"`
	Name          string     `gorm:"column:name" json:"name"`
	Description   *string    `gorm:"column:description" json:"description"`
	SlugCategory  *string    `gorm:"column:slug_category" json:"slug_category"`
	ImageCategory *string    `gorm:"column:image_category" json:"image_category"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Category) TableName() string { return "categories" }

// ---------------------------------------------------------------------------
// User Features
// ---------------------------------------------------------------------------

type ShippingAddress struct {
	ShippingAddressID int32      `gorm:"column:shipping_address_id;primaryKey" json:"shipping_address_id"`
	OrderID           int32      `gorm:"column:order_id" json:"order_id"`
	Alamat            string     `gorm:"column:alamat" json:"alamat"`
	Provinsi          string     `gorm:"column:provinsi" json:"provinsi"`
	Negara            string     `gorm:"column:negara" json:"negara"`
	Kota              string     `gorm:"column:kota" json:"kota"`
	Courier           string     `gorm:"column:courier" json:"courier"`
	ShippingMethod    string     `gorm:"column:shipping_method" json:"shipping_method"`
	ShippingCost      float64    `gorm:"column:shipping_cost" json:"shipping_cost"`
	CreatedAt         *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ShippingAddress) TableName() string { return "shipping_addresses" }

type Cart struct {
	CartID    int32      `gorm:"column:cart_id;primaryKey" json:"cart_id"`
	UserID    int32      `gorm:"column:user_id" json:"user_id"`
	ProductID int32      `gorm:"column:product_id" json:"product_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Price     int32      `gorm:"column:price" json:"price"`
	Image     string     `gorm:"column:image" json:"image"`
	Quantity  int32      `gorm:"column:quantity" json:"quantity"`
	Weight    int32      `gorm:"column:weight" json:"weight"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Cart) TableName() string { return "carts" }

// ---------------------------------------------------------------------------
// Reviews
// ---------------------------------------------------------------------------

type Review struct {
	ReviewID  int32      `gorm:"column:review_id;primaryKey" json:"review_id"`
	UserID    int32      `gorm:"column:user_id" json:"user_id"`
	ProductID int32      `gorm:"column:product_id" json:"product_id"`
	Name      string     `gorm:"column:name" json:"name"`
	Comment   string     `gorm:"column:comment" json:"comment"`
	Rating    int32      `gorm:"column:rating" json:"rating"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Review) TableName() string { return "reviews" }

type ReviewDetail struct {
	ReviewDetailID int32      `gorm:"column:review_detail_id;primaryKey" json:"review_detail_id"`
	ReviewID       int32      `gorm:"column:review_id" json:"review_id"`
	Type           string     `gorm:"column:type" json:"type"`
	Url            string     `gorm:"column:url" json:"url"`
	Caption        *string    `gorm:"column:caption" json:"caption"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (ReviewDetail) TableName() string { return "review_details" }

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

type Transaction struct {
	TransactionID int32      `gorm:"column:transaction_id;primaryKey" json:"transaction_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID    int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PaymentMethod string     `gorm:"column:payment_method" json:"payment_method"`
	Amount        int32      `gorm:"column:amount" json:"amount"`
	PaymentStatus string     `gorm:"column:payment_status" json:"payment_status"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Transaction) TableName() string { return "transactions" }

// ---------------------------------------------------------------------------
// Infrastructure
// ---------------------------------------------------------------------------

type OutboxEvent struct {
	OutboxID      int64      `gorm:"column:outbox_id;primaryKey" json:"outbox_id"`
	Topic         string     `gorm:"column:topic" json:"topic"`
	EventKey      string     `gorm:"column:event_key" json:"event_key"`
	Payload       []byte     `gorm:"column:payload;type:bytea" json:"payload"`
	Status        string     `gorm:"column:status" json:"status"`
	Attempts      int32      `gorm:"column:attempts" json:"attempts"`
	NextAttemptAt time.Time  `gorm:"column:next_attempt_at" json:"next_attempt_at"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (OutboxEvent) TableName() string { return "outbox_events" }

type ConsumerInbox struct {
	ConsumerName       string     `gorm:"column:consumer_name" json:"consumer_name"`
	EventKey           string     `gorm:"column:event_key" json:"event_key"`
	Topic              string     `gorm:"column:topic" json:"topic"`
	PartitionID        int32      `gorm:"column:partition_id" json:"partition_id"`
	MessageOffset      int64      `gorm:"column:message_offset" json:"message_offset"`
	Status             string     `gorm:"column:status" json:"status"`
	Attempts           int32      `gorm:"column:attempts" json:"attempts"`
	LeaseUntil         *time.Time `gorm:"column:lease_until" json:"lease_until"`
	LastError          string     `gorm:"column:last_error" json:"last_error"`
	ProcessedAt        *time.Time `gorm:"column:processed_at" json:"processed_at"`
	ReservationVersion int64      `gorm:"column:reservation_version" json:"reservation_version"`
}

func (ConsumerInbox) TableName() string { return "consumer_inbox" }
