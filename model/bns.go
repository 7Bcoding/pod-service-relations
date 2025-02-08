package model

import "time"

type Bns struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	ProductID uint32    `json:"product_id"`
	Cluster   string    `json:"cluster"`
	Azone     string    `json:"azone"`
	Idc       string    `json:"idc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
