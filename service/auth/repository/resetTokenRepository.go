package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/microservice-ecommerce-shared/errors"
	resettoken_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/reset_token_errors"
	"gorm.io/gorm"
)

type resetTokenRepository struct {
	db *gorm.DB
}

func NewResetTokenRepository(db *gorm.DB) *resetTokenRepository {
	return &resetTokenRepository{db: db}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*models.ResetToken, error) {
	var result models.ResetToken
	err := r.db.WithContext(ctx).Where("token = ?", code).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, resettoken_errors.ErrTokenNotFound.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &result, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}
	err = r.db.WithContext(ctx).Create(resetToken).Error
	if err != nil {
		return nil, resettoken_errors.ErrCreateResetToken.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) CreateResetTokenInTx(ctx context.Context, tx *gorm.DB, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, err
	}
	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}
	err = tx.WithContext(ctx).Create(resetToken).Error
	if err != nil {
		return nil, resettoken_errors.ErrCreateResetToken.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", user_id).Delete(&models.ResetToken{}).Error
	if err != nil {
		return resettoken_errors.ErrDeleteByUserID.WithInternal(err)
	}
	return nil
}
