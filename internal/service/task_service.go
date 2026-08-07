package service

import (
	"Tamrin/tasks/internal/domain"
	"Tamrin/tasks/internal/repository"
	"context"
	"errors"
)

type TaskServ struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskServ {
	return &TaskServ{repo: repo}
}

func (service *TaskServ) GetAll(ctx context.Context, userID int64) ([]domain.Task, error) {
	tasks, err := service.repo.GetAll(ctx, userID)
	return tasks, err
}

func (service *TaskServ) GetById(ctx context.Context, id int, userID int64) (domain.Task, error) {
	if id <= 0 {
		return domain.Task{}, errors.New("id must be positive")
	}
	return service.repo.GetById(ctx, id, userID)
}
func (service *TaskServ) Create(ctx context.Context, task domain.Task, userID int64) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("title is required")
	}

	id, err := service.repo.Create(ctx, task, userID)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("failed to create task")
	}
	return id, nil
}
func (service *TaskServ) Update(ctx context.Context, task domain.Task, id int, userID int64) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	if task.Title == "" {
		return errors.New("title is required")
	}
	return service.repo.Update(ctx, task, id, userID)
}
func (service *TaskServ) Delete(ctx context.Context, id int, userID int64) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	return service.repo.Delete(ctx, id, userID)
}
