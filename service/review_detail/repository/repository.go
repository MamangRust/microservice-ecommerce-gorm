package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
	ReviewDetailQuery   ReviewDetailQueryRepository
	ReviewDetailCommand ReviewDetailCommandRepository
}

func NewRepositories(DB *gorm.DB) *Repositories {
	return &Repositories{
		ReviewDetailQuery:   NewReviewDetailQueryRepository(DB),
		ReviewDetailCommand: NewReviewDetailCommandRepository(DB),
	}
}
