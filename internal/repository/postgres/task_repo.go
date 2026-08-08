package postgres

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
)

type RepoStruct struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *RepoStruct {
	return &RepoStruct{db: db}
}
func (r *RepoStruct) GetAll(ctx context.Context, userID int64) ([]domain.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, status, created_at FROM tasks WHERE user_id=$1", userID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *RepoStruct) GetById(ctx context.Context, id int, userID int64) (domain.Task, error) {
	var task domain.Task
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, description, status, created_at FROM tasks WHERE id=$1 AND user_id=$2",
		id, userID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (r *RepoStruct) Create(ctx context.Context, task domain.Task, userID int64) (int64, error) {
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO tasks (title, user_id, description, status) VALUES ($1, $2, $3, $4) RETURNING id",
		task.Title, userID, task.Description, task.Status).Scan(&task.ID)
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (r *RepoStruct) Update(ctx context.Context, task domain.Task, id int, userID int64) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4 AND user_id=$5",
		task.Title, task.Description, task.Status, id, userID)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *RepoStruct) Delete(ctx context.Context, id int, userID int64) error {
	rows, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	row, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return sql.ErrNoRows
	}
	return nil
}
