package model

import "time"

type Product struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Bnses     []*Bns    `json:"bnses,omitempty"`
}

type ProductService struct {
	ID        uint32    `json:"id"`
	ProductID uint32    `json:"product_id"`
	ServiceID uint32    `json:"service_id"`
	Enable    bool      `json:"enable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Product   *Product  `json:"product"`
	Service   *Service  `json:"service"`
}
