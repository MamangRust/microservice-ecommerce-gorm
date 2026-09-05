package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/microservice-ecommerce-shared/errors/merchant"
	"gorm.io/gorm"
)

type merchantDocumentCommandRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentCommandRepository(db *gorm.DB) *merchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{db: db}
}

func (r *merchantDocumentCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	doc := &models.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         stringPtr(""),
	}

	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return doc, nil
}

func (r *merchantDocumentCommandRepository) CreateInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	doc := &models.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         stringPtr(""),
	}

	if err := tx.WithContext(ctx).Create(doc).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	return doc, nil
}

func (r *merchantDocumentCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).First(&doc, *request.DocumentID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	updates := map[string]interface{}{
		"document_type": request.DocumentType,
		"document_url":  request.DocumentUrl,
		"status":        request.Status,
		"note":          stringPtr(request.Note),
	}

	if err := r.db.WithContext(ctx).Model(&doc).Updates(updates).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&doc, *request.DocumentID)
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).First(&doc, *request.DocumentID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	updates := map[string]interface{}{
		"status": request.Status,
		"note":   stringPtr(request.Note),
	}

	if err := r.db.WithContext(ctx).Model(&doc).Updates(updates).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	r.db.WithContext(ctx).First(&doc, *request.DocumentID)
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := tx.WithContext(ctx).First(&doc, *request.DocumentID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}

	updates := map[string]interface{}{
		"status": request.Status,
		"note":   stringPtr(request.Note),
	}

	if err := tx.WithContext(ctx).Model(&doc).Updates(updates).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}

	tx.WithContext(ctx).First(&doc, *request.DocumentID)
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Trash(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).First(&doc, documentID).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	if err := r.db.WithContext(ctx).Model(&doc).Update("deleted_at", time.Now()).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) Restore(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Unscoped().Where("document_id = ? AND deleted_at IS NOT NULL", documentID).First(&doc).Error; err != nil {
		return nil, merchant_errors.ErrMerchantNotFound
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&doc).Update("deleted_at", nil).Error; err != nil {
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}
	r.db.WithContext(ctx).Unscoped().First(&doc, documentID)
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("document_id = ?", documentID).Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return false, merchant_errors.ErrMerchantNotFound
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantDocument{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(result.Error)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, merchant_errors.ErrMerchantInternal.WithInternal(result.Error)
	}
	return true, nil
}
