package postgres

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	connStr := "user=postgres dbname=test_db password=12345678 sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic("can't connect to test database: " + err.Error())
	}
	defer testDB.Close()
	_, err = testDB.Exec(`
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT DEFAULT 'pending',
        created_at TIMESTAMP DEFAULT NOW()
    );
`)
	if err != nil {
		panic("can't create table: " + err.Error())
	}
	code := m.Run()
	testDB.Exec(`DROP TABLE tasks;`)
	os.Exit(code)
}
func cleanDB(t *testing.T) {
	_, err := testDB.Exec("DELETE FROM tasks")
	if err != nil {
		t.Errorf("can't clean database: %v", err)
	}
}
func TestCreateTask(t *testing.T) {
	cleanDB(t)
	repo := NewRepo(testDB)
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	ctx := context.Background()
	integer, err := repo.Create(ctx, task)
	if err != nil {
		t.Errorf("can't create task: %v", err)
		return
	}
	if integer == 0 {
		t.Errorf("integer should be greater than zero")
		return
	}
}
func TestGetTaskByID(t *testing.T) {
	cleanDB(t)
	repo := NewRepo(testDB)
	ctx := context.Background()
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	integer, err := repo.Create(ctx, task)
	if err != nil {
		t.Errorf("can't create task: %v", err)
		return
	}
	taskByID, err := repo.GetById(ctx, integer)
	if err != nil {
		t.Errorf("can't get task by id: %v", err)
		return
	}
	fmt.Println(taskByID)
}
func TestGetAllTask(t *testing.T) {
	cleanDB(t)
	repo := NewRepo(testDB)
	ctx := context.Background()
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	task2 := domain.Task{
		Title:       "test2 task",
		Description: "test2 description",
		Status:      "pending",
	}
	task3 := domain.Task{
		Title:       "test3 task",
		Description: "test3 description",
		Status:      "pending",
	}
	integer, err := repo.Create(ctx, task)
	if err != nil {
		t.Errorf("can't create task: %v", err)
		return
	}
	if integer == 0 {
		t.Errorf("integer should be greater than zero")
		return
	}
	integer, err = repo.Create(ctx, task2)
	if err != nil {
		t.Errorf("can't create task: %v", err)
		return
	}
	if integer == 0 {
		t.Errorf("integer should be greater than zero")
		return
	}
	integer, err = repo.Create(ctx, task3)
	if err != nil {
		t.Errorf("can't create task: %v", err)
		return
	}
	if integer == 0 {
		t.Errorf("integer should be greater than zero")
		return
	}
	tasks, err := repo.GetAll(ctx)
	if err != nil {
		t.Errorf("can't get all tasks: %v", err)
		return
	}
	if len(tasks) == 0 {
		t.Errorf("no tasks found")
		return
	}
	fmt.Println(tasks)
}
func TestUpdateTask(t *testing.T) {
	cleanDB(t)
	repo := NewRepo(testDB)
	ctx := context.Background()

	// ۱. یک تسک بساز
	task := domain.Task{
		Title:       "before update",
		Description: "before description",
		Status:      "pending",
	}
	id, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("can't create task: %v", err)
	}

	// ۲. تغییرات رو اعمال کن
	task.Title = "after update"
	task.Description = "after description"
	task.Status = "done"

	err = repo.Update(ctx, task, id)
	if err != nil {
		t.Fatalf("can't update task: %v", err)
	}

	// ۳. تغییرات رو چک کن
	updated, err := repo.GetById(ctx, id)
	if err != nil {
		t.Fatalf("can't get task: %v", err)
	}

	if updated.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, updated.Title)
	}
	if updated.Description != task.Description {
		t.Errorf("expected description %q, got %q", task.Description, updated.Description)
	}
	if updated.Status != task.Status {
		t.Errorf("expected status %q, got %q", task.Status, updated.Status)
	}
}
func TestDeleteTask(t *testing.T) {
	cleanDB(t)
	repo := NewRepo(testDB)
	ctx := context.Background()
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	id, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("can't create task: %v", err)
		return
	}
	fmt.Println(id)
	err = repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("can't delete task: %v", err)
		return
	}
	fmt.Println("task deleted")
}
