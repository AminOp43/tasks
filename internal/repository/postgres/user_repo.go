package postgres

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type UserRepoStruct struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepoStruct {
	return &UserRepoStruct{
		db: db,
	}
}
func (u *UserRepoStruct) Create(ctx context.Context, user domain.User) (int64, error) {
	err := u.db.QueryRowContext(
		ctx,
		`INSERT INTO users (username, password_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		user.Username,
		user.PasswordHash,
	).Scan(&user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}
func (u *UserRepoStruct) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := u.db.QueryRowContext(ctx, `SELECT  id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`,
		username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}
	return &user, nil
}
