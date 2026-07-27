package service

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"testing"
)

type MockRepo struct {
	GetAllFunc  func(ctx context.Context) ([]domain.Task, error)
	GetByIdFunc func(ctx context.Context, id int) (domain.Task, error)
	CreateFunc  func(ctx context.Context, task domain.Task) (int, error)
	UpdateFunc  func(ctx context.Context, task domain.Task, id int) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *MockRepo) GetAll(ctx context.Context) ([]domain.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepo) GetById(ctx context.Context, id int) (domain.Task, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return domain.Task{}, nil
}

func (m *MockRepo) Create(ctx context.Context, task domain.Task) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return 0, nil
}

func (m *MockRepo) Update(ctx context.Context, task domain.Task, id int) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task, id)
	}
	return nil
}

func (m *MockRepo) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
func TestMockGetAll(t *testing.T) {
	mockRepo := &MockRepo{
		GetAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return []domain.Task{{Title: "title test", Description: "desc test", Status: "pending"}}, nil
		},
	}
	newTaskService := NewTaskService(mockRepo)
	tasks, err := newTaskService.GetAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "title test" {
		t.Errorf("expected title 'title test', got '%s'", tasks[0].Title)
	}
	t.Log(tasks)
}
func TestMockGetById(t *testing.T) {
	mockRepo := &MockRepo{
		GetByIdFunc: func(ctx context.Context, id int) (domain.Task, error) {
			return domain.Task{
				ID:          1,
				Title:       "title test",
				Description: "desc test",
				Status:      "pending",
			}, nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task, err := newTaskSrv.GetById(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	if task.Title == "" {
		t.Errorf("expected non-empty title, got empty")
	}
	if task.ID == 0 {
		t.Errorf("expected non-zero ID, got %d", task.ID)
	}
	t.Log(task)
}
func TestMockCreate(t *testing.T) {
	mockRepo := &MockRepo{
		CreateFunc: func(ctx context.Context, task domain.Task) (int, error) {
			return 1, nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "title test",
		Description: "desc test",
		Status:      "pending",
	}
	id, err := newTaskSrv.Create(context.Background(), task)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Errorf("expected 1, got %d", id)
	}
	t.Log(task)
}
func TestMockUpdate(t *testing.T) {
	mockRepo := &MockRepo{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int) error {
			return nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "title test",
		Description: "desc test",
		Status:      "pending",
	}
	err := newTaskSrv.Update(context.Background(), task, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("succesfully updated : %+v", task)
}
func TestMockDelete(t *testing.T) {
	mockRepo := &MockRepo{
		DeleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}
	newTaskSrv := NewTaskService(mockRepo)
	err := newTaskSrv.Delete(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("successfully deleted")
}
func TestMockGetById_NotFound(t *testing.T) {
	mockRepo := &MockRepo{
		GetByIdFunc: func(ctx context.Context, id int) (domain.Task, error) {
			return domain.Task{}, sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	_, err := svc.GetById(context.Background(), 999)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	t.Log(" GetById not found test passed!")
}
func TestMockUpdate_NotFound(t *testing.T) {
	mockRepo := &MockRepo{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int) error {
			return sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	err := svc.Update(context.Background(), task, 999)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	t.Log(" Update not found test passed!")
}
func TestMockDelete_NotFound(t *testing.T) {
	mockRepo := &MockRepo{
		DeleteFunc: func(ctx context.Context, id int) error {
			return sql.ErrNoRows
		},
	}
	svc := NewTaskService(mockRepo)
	err := svc.Delete(context.Background(), 999)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	t.Log(" Delete not found test passed!")
}
