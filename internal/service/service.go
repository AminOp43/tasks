package service

import (
	"Tamrin/tasks/internal/domain"
	"context"
)

type IntService interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetById(ctx context.Context, id int) (domain.Task, error)
	Create(ctx context.Context, task domain.Task) (int, error)
	Update(ctx context.Context, task domain.Task, id int) error
	Delete(ctx context.Context, id int) error
}
