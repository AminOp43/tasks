package domain

import "time"

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"text"`
	Description string    `json:"desc"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"creates_at"`
}
