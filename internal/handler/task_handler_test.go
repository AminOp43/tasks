package handler

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	GetAllFunc  func(ctx context.Context) ([]domain.Task, error)
	GetByIdFunc func(ctx context.Context, id int) (domain.Task, error)
	CreateFunc  func(ctx context.Context, task domain.Task) (int, error)
	UpdateFunc  func(ctx context.Context, task domain.Task, id int) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *MockService) GetAll(ctx context.Context) ([]domain.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}
func (m *MockService) GetById(ctx context.Context, id int) (domain.Task, error) {
	if m.GetByIdFunc != nil {
		return m.GetByIdFunc(ctx, id)
	}
	return domain.Task{}, nil
}

func (m *MockService) Create(ctx context.Context, task domain.Task) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return 0, nil
}

func (m *MockService) Update(ctx context.Context, task domain.Task, id int) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task, id)
	}
	return nil
}

func (m *MockService) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
func TestHandler_GetTasks_Success(t *testing.T) {
	mockService := &MockService{
		GetAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return []domain.Task{}, nil
		},
	}
	tasks, err := mockService.GetAll(context.Background())
	if err != nil {
		t.Error(err)
	}
	tHandler := NewTaskHandler(mockService)
	r := gin.Default()
	r.GET("/tasks", tHandler.GetAll)
	//یک درخواست HTTP ساختگی می‌سازه
	req := httptest.NewRequest("GET", "/tasks", nil)

	//یک پاسخ‌گیر (Response Recorder) می‌سازه که پاسخ رو ضبط میکنه
	//جزئیات:
	//    httptest.NewRecorder() یک ResponseWriter ساختگی هست که پاسخ رو توی حافظه ذخیره میکنه.
	//
	// بعداً می‌تونی w.Code و w.Body.String() رو چک کنی
	w := httptest.NewRecorder()

	//درخواست رو به gin.Engine می‌ده تا پردازش کنه
	//
	//جزئیات:
	//
	//    r.ServeHTTP دقیقاً مثل اینه که سرور روشن باشه و یک درخواست واقعی بهش برسه.
	//
	//   این خط Handler رو اجرا میکنه و پاسخ رو توی w ذخیره میکنه
	r.ServeHTTP(w, req)

	//w.Code کد وضعیت HTTP رو برمی‌گردونه.
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	t.Log(tasks)
}
func TestHandler_GetTasks_Error(t *testing.T) {
	mockService := &MockService{
		GetAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return nil, errors.New("test error")
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
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
		GetByIdFunc: func(ctx context.Context, id int) (domain.Task, error) {
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
		GetByIdFunc: func(ctx context.Context, id int) (domain.Task, error) {
			return domain.Task{}, sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.GET("/tasks/:id", handler.GetById)
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
func TestHandler_GetTask_InvalidID(t *testing.T) {
	// نیازی به MockService نیست چون Handler قبل از رسیدن به Service، ID رو چک میکنه
	handler := NewTaskHandler(&MockService{}) // ← می‌تونی خالی بدی

	r := gin.Default()
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
		CreateFunc: func(ctx context.Context, task domain.Task) (int, error) {
			return 1, nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.POST("/tasks", handler.Create)

	jsonData := `{"title":"test task","description":"test description","status":"pending"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// چک کردن بدنه‌ی پاسخ
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
		UpdateFunc: func(ctx context.Context, task domain.Task, id int) error {
			return nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.PUT("/tasks/:id", handler.Update)
	jsonData := `{"title":"test task","description":"test description","status":"pending"}`
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
func TestHandler_UpdateTask_NotFound(t *testing.T) {
	mockService := &MockService{
		UpdateFunc: func(ctx context.Context, task domain.Task, id int) error {
			return sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.PUT("/tasks/:id", handler.Update)
	jsonData := `{"title":"test task","description":"test description","status":"pending"}`
	req := httptest.NewRequest("PUT", "/tasks/999", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
func TestHandler_DeleteTask_Success(t *testing.T) {
	mockService := &MockService{
		DeleteFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.DELETE("/tasks/:id", handler.Delete)
	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
func TestHandler_DeleteTask_NotFound(t *testing.T) {
	mockService := &MockService{
		DeleteFunc: func(ctx context.Context, id int) error {
			return sql.ErrNoRows
		},
	}
	handler := NewTaskHandler(mockService)
	r := gin.Default()
	r.DELETE("/tasks/:id", handler.Delete)
	req := httptest.NewRequest("DELETE", "/tasks/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
