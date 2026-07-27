package postgres

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"errors"
	"time"
)

type RepoStruct struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *RepoStruct {
	return &RepoStruct{db: db}
}
func (r *RepoStruct) GetAll(ctx context.Context) ([]domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, status, created_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		scanErr := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}
		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *RepoStruct) GetById(ctx context.Context, id int) (domain.Task, error) {
	var task domain.Task
	err := r.db.QueryRowContext(ctx, "SELECT id, title, description, status, created_at FROM tasks WHERE id=$1",
		id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, err
		}
		return task, err
	}
	return task, nil
}

func (r *RepoStruct) Create(ctx context.Context, task domain.Task) (int, error) {
	task.CreatedAt = time.Now()
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO tasks (title, description, status,created_at) VALUES ($1, $2, $3,$4) RETURNING id",
		task.Title, task.Description, task.Status, task.CreatedAt).Scan(&task.ID)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (r *RepoStruct) Update(ctx context.Context, task domain.Task, id int) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
		task.Title, task.Description, task.Status, id)
	if err != nil {
		return err
	}
	row, _ := res.RowsAffected()
	if row == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RepoStruct) Delete(ctx context.Context, id int) error {
	rows, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	row, _ := rows.RowsAffected()
	if row == 0 {
		return sql.ErrNoRows
	}
	return nil
}
