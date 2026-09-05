package repository

import "gorm.io/gorm"

type Repositories struct {
	SliderQuery   SliderQueryRepository
	SliderCommand SliderCommandRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		SliderQuery:   NewSliderQueryRepository(db),
		SliderCommand: NewSliderCommandRepository(db),
	}
}
