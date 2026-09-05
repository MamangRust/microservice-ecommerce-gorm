package repository

import "gorm.io/gorm"

type Repositories struct {
	ShippingAddressQuery   ShippingAddressQueryRepository
	ShippingAddressCommand ShippingAddressCommandRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		ShippingAddressQuery:   NewShippingAddressQueryRepository(db),
		ShippingAddressCommand: NewShippingAddressCommandRepository(db),
	}
}
