package handler

import (
	"Tamrin/tasks/internal/domain"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	GetAllFunc  func(ctx context.Context, userID int64) ([]domain.Task, error)
	GetByIdFunc func(ctx context.Context, id int, userID int64) (domain.Task, error)
	CreateFunc  func(ctx context.Context, task domain.Task, userID int64) (int64, error)
	UpdateFunc  func(ctx context.Context, task domain.Task, id int, userID int64) error
	DeleteFunc  func(ctx context.Context, id int, userID int64) error
}

func (m *MockService) GetAll(ctx context.Context, userID int64) ([]domain.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, userID)
	}
	return nil, nil
}
func (m *MockService) GetById(ctx context.Context, id int, userID int64) (domain.Task, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id, userID)
	}
	return domain.Task{}, nil
}

func (m *MockService) Create(ctx context.Context, task domain.Task, userID int64) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task, userID)
	}
	return 0, nil
}

func (m *MockService) Update(ctx context.Context, task domain.Task, id int, userID int64) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task, id, userID)
	}
	return nil
}

func (m *MockService) Delete(ctx context.Context, id int, userID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, userID)
	}
	return nil
}
func setTestUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	}
}

func TestHandler_GetTasks_Success(t *testing.T) {
	mockService := &MockService{
		GetAllFunc: func(ctx context.Context, userID int64) ([]domain.Task, error) {
			if userID != 1 {
				t.Errorf("expected userID 1, got %d", userID)
			}
			return []domain.Task{{
				ID:          1,
				Title:       "test task",
				Description: "test description",
				Status:      "pending",
			}}, nil
		},
	}
	tHandler := NewTaskHandler(mockService)

	r := gin.Default()
	r.Use(setTestUserID())
	r.GET("/tasks", tHandler.GetAll)

	req := httptest.NewRequest("GET", "/tasks", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var response struct {
		Tasks []domain.Task `json:"tasks"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("cannot parse response: %v", err)
	}

	if len(response.Tasks) != 1 {
		t.Fatalf(
			"expected 1 task, got %d",
			len(response.Tasks),
		)
	}

	task := response.Tasks[0]

	if task.ID != 1 {
		t.Errorf("expected task ID 1, got %d", task.ID)
	}

	if task.Title != "test task" {
		t.Errorf(
			"expected title %q, got %q",
			"test task",
			task.Title,
		)
	}
}
func TestHandler_GetTasks_Error(t *testing.T) {
	mockService := &MockService{
		GetAllFunc: func(ctx context.Context, userID int64) ([]domain.Task, error) {
			return nil, errors.New("test error")
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.GET("/tasks", handler.GetAll)
	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
func TestHandler_GetTask_Success(t *testing.T) {
	mockService := &MockService{
		GetByIdFunc: func(ctx context.Context, id int, userID int64) (domain.Task, error) {
			if id != 1 {
				t.Errorf("expected task ID 1, got %d", id)
			}

			if userID != 1 {
				t.Errorf("expected user ID 1, got %d", userID)
			}
			return domain.Task{
				ID:          1,
				Title:       "test task",
				Description: "test description",
				Status:      "pending",
			}, nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.GET("/tasks/:id", handler.GetById)

	req := httptest.NewRequest("GET", "/tasks/1", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	var response map[string]domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("can't parse response: %v", err)
	}

	task := response["task"]
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Title != "test task" {
		t.Errorf("expected title 'test task', got '%s'", task.Title)
	}
}
func TestHandler_GetTask_NotFound(t *testing.T) {
	mockService := &MockService{
		GetByIdFunc: func(ctx context.Context, id int, userID int64) (domain.Task, error) {
			return domain.Task{}, sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.GET("/tasks/:id", handler.GetById)
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
func TestHandler_GetTask_InvalidID(t *testing.T) {
	handler := NewTaskHandler(&MockService{})
	r := gin.Default()
	r.Use(setTestUserID())
	r.GET("/tasks/:id", handler.GetById)
	req := httptest.NewRequest("GET", "/tasks/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
func TestHandler_CreateTask_Success(t *testing.T) {
	mockService := &MockService{
		CreateFunc: func(ctx context.Context, task domain.Task, userID int64) (int64, error) {
			if userID != 1 {
				t.Errorf("expected userID 1, got %d", userID)
			}
			return 1, nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.POST("/tasks", handler.Create)

	requestTask := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}

	jsonData, err := json.Marshal(requestTask)
	if err != nil {
		t.Fatalf("can't create request JSON: %v", err)
	}
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var response map[string]domain.Task
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("can't parse response: %v", err)
	}

	task := response["task"]
	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Title != "test task" {
		t.Errorf("expected title 'test task', got '%s'", task.Title)
	}
}
func TestHandler_UpdateTask_Success(t *testing.T) {
	mockService := &MockService{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int, userID int64) error {
			return nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.PUT("/tasks/:id", handler.Update)
	requestTask := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}

	jsonData, err := json.Marshal(requestTask)
	if err != nil {
		t.Fatalf("can't create request JSON: %v", err)
	}
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
func TestHandler_UpdateTask_NotFound(t *testing.T) {
	mockService := &MockService{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int, userID int64) error {
			return sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.PUT("/tasks/:id", handler.Update)
	requestTask := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}

	jsonData, err := json.Marshal(requestTask)
	if err != nil {
		t.Fatalf("can't create request JSON: %v", err)
	}
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
func TestHandler_DeleteTask_Success(t *testing.T) {
	mockService := &MockService{
		DeleteFunc: func(ctx context.Context, id int, userID int64) error {
			return nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.DELETE("/tasks/:id", handler.Delete)
	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
func TestHandler_DeleteTask_NotFound(t *testing.T) {
	mockService := &MockService{
		DeleteFunc: func(ctx context.Context, id int, userID int64) error {
			return sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.Use(setTestUserID())
	r.DELETE("/tasks/:id", handler.Delete)
	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
