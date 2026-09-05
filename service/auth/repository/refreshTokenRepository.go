package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	refreshtoken_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/refresh_token_errors"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *refreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var result models.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error) {
	var result models.RefreshToken
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, refreshtoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	refreshToken := &models.RefreshToken{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	}
	err = r.db.WithContext(ctx).Create(refreshToken).Error
	if err != nil {
		return nil, refreshtoken_errors.ErrCreateRefreshToken.WithInternal(err)
	}
	return refreshToken, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, refreshtoken_errors.ErrParseDate
	}

	var result models.RefreshToken
	err = r.db.WithContext(ctx).Where("user_id = ?", req.UserId).First(&result).Error
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}

	result.Token = req.Token
	result.Expiration = expirationTime
	err = r.db.WithContext(ctx).Save(&result).Error
	if err != nil {
		return nil, refreshtoken_errors.ErrUpdateRefreshToken.WithInternal(err)
	}
	return &result, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return refreshtoken_errors.ErrDeleteRefreshToken.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).Delete(&models.RefreshToken{}).Error
	if err != nil {
		return refreshtoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
