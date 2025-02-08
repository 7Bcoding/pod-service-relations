package model

import "time"

type User struct {
	ID        uint32    `json:"id"`
	Uid       string    `json:"uid"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
