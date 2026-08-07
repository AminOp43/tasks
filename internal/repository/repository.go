package repository

import (
	"Tamrin/tasks/internal/domain"
	"context"
)

type TaskRepository interface {
	GetAll(ctx context.Context, userID int64) ([]domain.Task, error)
	GetById(ctx context.Context, id int, userID int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task, userID int64) (int64, error)
	Update(ctx context.Context, task domain.Task, id int, userID int64) error
	Delete(ctx context.Context, id int, userID int64) error
}
type UserRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}
