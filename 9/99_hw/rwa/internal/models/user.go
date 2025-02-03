package models

import "time"

type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
	Bio       string    `json:"bio"`
	Image     string    `json:"image"`
	Follows   map[uint64]struct{}
}
