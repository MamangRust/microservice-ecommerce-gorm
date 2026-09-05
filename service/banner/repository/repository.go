package repository

import "gorm.io/gorm"

type Repositories struct {
	BannerQuery   BannerQueryRepository
	BannerCommand BannerCommandRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		BannerQuery:   NewBannerQueryRepository(db),
		BannerCommand: NewBannerCommandRepository(db),
	}
}
