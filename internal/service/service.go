package service

import (
	"Tamrin/tasks/internal/domain"
	"context"
)

type TaskService interface {
	GetAll(ctx context.Context, userID int64) ([]domain.Task, error)
	GetById(ctx context.Context, id int, userID int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task, userID int64) (int64, error)
	Update(ctx context.Context, task domain.Task, id int, userID int64) error
	Delete(ctx context.Context, id int, userID int64) error
}
type UserService interface {
	SignUp(ctx context.Context, req domain.AuthRequest) error
	Login(ctx context.Context, req domain.AuthRequest) (string, error)
}
