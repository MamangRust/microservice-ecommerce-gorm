package repository

import (
	"context"
	"time"

	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

type MerchantCertResult struct {
	MerchantCertificationID int32
	MerchantID              int32
	Title                   string
	Description             *string
	IssuedBy                *string
	IssueDate               *time.Time
	ExpiryDate              *time.Time
	CertificateUrl          *string
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
	DeletedAt               *time.Time
	TotalCount              int64
}

type MerchantAwardQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantCertResult, error)
	FindByID(ctx context.Context, id int) (*MerchantCertResult, error)
}

type MerchantAwardCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error)
	Update(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*models.MerchantCertificationsAndAward, error)
	Trash(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error)
	Restore(ctx context.Context, id int) (*models.MerchantCertificationsAndAward, error)
	DeletePermanent(ctx context.Context, id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// MerchantQueryRepository is the shared merchant query contract, provided by the
// merchant gRPC adapter so this service never holds a raw gRPC client.
type MerchantQueryRepository = merchantadapter.QueryRepository
