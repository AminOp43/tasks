package domain

import "time"

type User struct {
	ID           int64     `json:"-"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
