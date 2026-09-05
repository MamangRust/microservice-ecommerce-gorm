package repository

import (
	pbuser "github.com/MamangRust/microservice-ecommerce-grpc/pb/user"
	"gorm.io/gorm"
)

type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

func NewRepositories(DB *gorm.DB, userQuery pbuser.UserQueryServiceClient) *Repositories {
	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(DB),
		MerchantCommand:         NewMerchantCommandRepository(DB),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(DB),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(DB),
		UserQuery:               NewUserQueryRepository(userQuery),
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
