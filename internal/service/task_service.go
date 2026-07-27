package service

import (
	"Tamrin/tasks/internal/domain"
	"Tamrin/tasks/internal/repository"
	"context"
	"errors"
)

type TaskService struct {
	repo repository.RepoInt
}

func NewTaskService(repo repository.RepoInt) *TaskService {
	return &TaskService{repo: repo}
}
func (service *TaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	tasks, err := service.repo.GetAll(ctx)
	return tasks, err
}

func (service *TaskService) GetById(ctx context.Context, id int) (domain.Task, error) {
	if id <= 0 {
		return domain.Task{}, errors.New("id must be positive")
	}
	return service.repo.GetById(ctx, id)
}
func (service *TaskService) Create(ctx context.Context, task domain.Task) (int, error) {
	if task.Title == "" {
		return 0, errors.New("title is required")
	}
	ser, err := service.repo.Create(ctx, task)
	if ser == 0 {
		return 0, errors.New("failed to create task")
	}
	return ser, err
}
func (service *TaskService) Update(ctx context.Context, task domain.Task, id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	if task.Title == "" {
		return errors.New("title is required")
	}
	return service.repo.Update(ctx, task, id)
}
func (service *TaskService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	return service.repo.Delete(ctx, id)
}
