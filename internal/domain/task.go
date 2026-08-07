package domain

import "time"

type Task struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"_"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
