package repository

import (
	"context"

	merchantadapter "github.com/MamangRust/microservice-ecommerce-pkg/adapter/merchant"
	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
)

// MerchantQueryRepository is the shared merchant query contract, provided by the
// merchant gRPC adapter so this service never holds a raw gRPC client.
type MerchantQueryRepository = merchantadapter.QueryRepository

type MerchantDetailResult struct {
	MerchantDetailID int32
	MerchantID       int32
	DisplayName      *string
	CoverImageUrl    *string
	LogoUrl          *string
	ShortDescription *string
	WebsiteUrl       *string
	CreatedAt        *string
	UpdatedAt        *string
	DeletedAt        *string
	TotalCount       int64
}

type MerchantDetailQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*MerchantDetailResult, error)
	FindByID(ctx context.Context, user_id int) (*MerchantDetailResult, error)
}

type MerchantDetailCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*models.MerchantDetail, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*models.MerchantDetail, error)
	Trash(ctx context.Context, merchant_id int) (*models.MerchantDetail, error)
	Restore(ctx context.Context, merchant_id int) (*models.MerchantDetail, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantSocialLinkCommandRepository interface {
	Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error)
	Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*models.MerchantSocialMediaLink, error)
	Trash(ctx context.Context, socialID int) (bool, error)
	Restore(ctx context.Context, socialID int) (bool, error)
	DeletePermanent(ctx context.Context, socialID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
