package postgres

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"regexp"
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock database: %v", err)
	}
	defer db.Close()
	repo := NewRepo(db)
	task := domain.Task{
		Title:       "test task",
		Description: "test description",
		Status:      "pending",
	}
	query := `
		INSERT INTO tasks
			(title, user_id, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(
			task.Title,
			int64(1),
			task.Description,
			task.Status,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(1),
		)
	id, err := repo.Create(context.Background(), task, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
func TestGetTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock database: %v", err)
	}
	defer db.Close()
	repo := NewRepo(db)
	query := `
		SELECT id, title, description, status, created_at 
		FROM tasks 
		WHERE id=$1 AND user_id=$2
	`
	createdAt := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1, int64(1)).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id",
					"title",
					"description",
					"status",
					"created_at",
				},
			).AddRow(
				1,
				"test task",
				"test description",
				"pending",
				createdAt,
			),
		)
	task, err := repo.GetById(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("can't get task by id: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Title != "test task" {
		t.Errorf(
			"expected title %q, got %q",
			"test task",
			task.Title,
		)
	}
	if task.Description != "test description" {
		t.Errorf(
			"expected description %q, got %q",
			"test description",
			task.Description,
		)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
func TestGetAllTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock database: %v", err)
	}
	defer db.Close()
	repo := NewRepo(db)
	createdAt := time.Now()
	query := `SELECT id, title, description, status, created_at FROM tasks WHERE user_id=$1`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(int64(1)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at"}).
				AddRow(1, "test task", "test description", "pending", createdAt),
		)
	tasks, err := repo.GetAll(context.Background(), 1)
	if err != nil {
		t.Fatalf("can't get all tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	if tasks[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", tasks[0].ID)
	}

	if tasks[0].Title != "test task" {
		t.Errorf(
			"expected title %q, got %q",
			"test task",
			tasks[0].Title,
		)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock database: %v", err)
	}
	defer db.Close()
	repo := NewRepo(db)
	task := domain.Task{
		Title:       "test updated task",
		Description: "test updated description",
		Status:      "updated pending"}
	query := `UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4 AND user_id=$5`
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(
			task.Title,
			task.Description,
			task.Status,
			1,
			int64(1),
		).WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.Update(context.Background(), task, 1, 1)
	if err != nil {
		t.Fatalf("can't update task: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("can't create mock database: %v", err)
	}
	defer db.Close()
	repo := NewRepo(db)
	query := `DELETE FROM tasks WHERE id=$1 AND user_id=$2`
	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(1, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	err = repo.Delete(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("can't delete task: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("database expectations were not met: %v", err)
	}
}
