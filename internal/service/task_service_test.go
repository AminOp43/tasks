package service

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"errors"
	"testing"
)

type MockTaskRepo struct {
	GetAllFunc  func(ctx context.Context, userID int64) ([]domain.Task, error)
	GetByIdFunc func(ctx context.Context, id int, userID int64) (domain.Task, error)
	CreateFunc  func(ctx context.Context, task domain.Task, userID int64) (int64, error)
	UpdateFunc  func(ctx context.Context, task domain.Task, id int, userID int64) error
	DeleteFunc  func(ctx context.Context, id int, userID int64) error
}

func (m *MockTaskRepo) GetAll(ctx context.Context, userID int64) ([]domain.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockTaskRepo) GetById(ctx context.Context, id int, userID int64) (domain.Task, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id, userID)
	}
	return domain.Task{}, nil
}

func (m *MockTaskRepo) Create(ctx context.Context, task domain.Task, userID int64) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task, userID)
	}
	return 0, nil
}

func (m *MockTaskRepo) Update(ctx context.Context, task domain.Task, id int, userID int64) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task, id, userID)
	}
	return nil
}

func (m *MockTaskRepo) Delete(ctx context.Context, id int, userID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, userID)
	}
	return nil
}
func TestTaskService_GetAll_Success(t *testing.T) {
	mockRepo := &MockTaskRepo{
		GetAllFunc: func(ctx context.Context, userID int64) ([]domain.Task, error) {
			return []domain.Task{{Title: "title test", Description: "desc test", Status: "pending"}}, nil
		},
	}
	newTaskService := NewTaskService(mockRepo)
	tasks, err := newTaskService.GetAll(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "title test" {
		t.Errorf("expected title 'title test', got '%s'", tasks[0].Title)
	}
	t.Log(tasks)
}
func TestTaskService_GetByID_Success(t *testing.T) {
	mockRepo := &MockTaskRepo{
		GetByIdFunc: func(ctx context.Context, id int, userID int64) (domain.Task, error) {
			return domain.Task{
				ID:          1,
				Title:       "title test",
				Description: "desc test",
				Status:      "pending",
			}, nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task, err := newTaskSrv.GetById(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if task.Title != "title test" {
		t.Errorf("expected title %q, got %q", "title test", task.Title)
	}
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	t.Log(task)
}
func TestTaskService_Create_Success(t *testing.T) {
	mockRepo := &MockTaskRepo{
		CreateFunc: func(ctx context.Context, task domain.Task, userID int64) (int64, error) {
			return 1, nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "title test",
		Description: "desc test",
		Status:      "pending",
	}
	id, err := newTaskSrv.Create(context.Background(), task, 1)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Errorf("expected 1, got %d", id)
	}
	t.Log(task)
}
func TestTaskService_Update_Success(t *testing.T) {
	mockRepo := &MockTaskRepo{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int, userID int64) error {
			return nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "title test",
		Description: "desc test",
		Status:      "pending",
	}
	err := newTaskSrv.Update(context.Background(), task, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("succesfully updated : %+v", task)
}
func TestTaskService_Delete_Success(t *testing.T) {
	mockRepo := &MockTaskRepo{
		DeleteFunc: func(ctx context.Context, id int, userID int64) error {
			return nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	err := newTaskSrv.Delete(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("successfully deleted")
}
func TestTaskService_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepo{
		GetByIdFunc: func(ctx context.Context, id int, userID int64) (domain.Task, error) {
			return domain.Task{}, sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	_, err := svc.GetById(context.Background(), 999, 1)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
	t.Log(" GetById not found test passed!")
}
func TestTaskService_Update_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepo{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int, userID int64) error {
			return sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	err := svc.Update(context.Background(), task, 999, 1)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
	t.Log(" Update not found test passed!")
}
func TestTaskService_Delete_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepo{
		DeleteFunc: func(ctx context.Context, id int, userID int64) error {
			return sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	err := svc.Delete(context.Background(), 999, 1)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
	t.Log(" Delete not found test passed!")
}
