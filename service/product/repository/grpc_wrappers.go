package repository

import (
	pbcategory "github.com/MamangRust/microservice-ecommerce-grpc/pb/category"
	pbmerchant "github.com/MamangRust/microservice-ecommerce-grpc/pb/merchant"
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/errors/category_errors"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
)

func (r *categoryQueryRepository) FindByID(ctx context.Context, category_id int) (*models.Category, error) {
	res, err := r.client.FindById(ctx, &pbcategory.FindByIdCategoryRequest{Id: int32(category_id)})
	if err != nil {
		return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
	}

	data := res.Data
	return &models.Category{
		CategoryID:    data.Id,
		Name:          data.Name,
		Description:   &data.Description,
		SlugCategory:  &data.SlugCategory,
		ImageCategory: &data.ImageCategory,
		CreatedAt:     parseTimeString(data.CreatedAt),
		UpdatedAt:     parseTimeString(data.UpdatedAt),
	}, nil
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, user_id int) (*models.Merchant, error) {
	res, err := r.client.FindById(ctx, &pbmerchant.FindByIdMerchantRequest{Id: int32(user_id)})
	if err != nil {
		return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
	}

	data := res.Data
	return &models.Merchant{
		MerchantID:   data.Id,
		UserID:       data.UserId,
		Name:         data.Name,
		Description:  &data.Description,
		Address:      &data.Address,
		ContactEmail: &data.ContactEmail,
		ContactPhone: &data.ContactPhone,
		Status:       data.Status,
		CreatedAt:    parseTimeString(data.CreatedAt),
		UpdatedAt:    parseTimeString(data.UpdatedAt),
	}, nil
}

func parseTimeString(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05.000", s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", s)
		if err != nil {
			return nil
		}
	}
	return &t
}
