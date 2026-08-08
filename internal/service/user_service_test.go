package service

import (
	"Tamrin/tasks/internal/domain"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type MockUserRepo struct {
	CreateFunc         func(ctx context.Context, user domain.User) (int64, error)
	FindByUsernameFunc func(ctx context.Context, username string) (*domain.User, error)
}

func (m *MockUserRepo) Create(ctx context.Context, user domain.User) (int64, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return 0, nil
}
func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	if m.FindByUsernameFunc != nil {
		return m.FindByUsernameFunc(ctx, username)
	}

	return nil, nil
}
func TestUserServ_SignUp(t *testing.T) {
	repo := &MockUserRepo{
		CreateFunc: func(ctx context.Context, user domain.User) (int64, error) {
			if user.Username != "test" {
				t.Errorf(
					"expected username %q, got %q",
					"test",
					user.Username,
				)
			}
			if user.PasswordHash == "" {
				t.Error("expected password hash, got empty")
			}
			if user.PasswordHash == "test" {
				t.Error("password should not be stored as plain text")
			}
			return 1, nil
		},
	}
	userService := NewUserServ(repo)
	user := domain.AuthRequest{
		Username: "test",
		Password: "test",
	}
	err := userService.SignUp(context.Background(), user)
	if err != nil {
		t.Error(err)
	}
}
func TestUserServ_SignUp_DuplicateUsername(t *testing.T) {
	errDuplicateUsername := errors.New("username already exists")

	repo := &MockUserRepo{
		CreateFunc: func(
			ctx context.Context,
			user domain.User,
		) (int64, error) {
			return 0, errDuplicateUsername
		},
	}

	userService := NewUserServ(repo)

	req := domain.AuthRequest{
		Username: "test",
		Password: "test",
	}

	err := userService.SignUp(context.Background(), req)

	if !errors.Is(err, errDuplicateUsername) {
		t.Errorf(
			"expected duplicate username error, got %v",
			err,
		)
	}
}
func TestUserServ_Login(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("test"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("can't hash password: %v", err)
	}

	t.Setenv("JWT_SECRET", "test-secret")

	repo := &MockUserRepo{
		FindByUsernameFunc: func(
			ctx context.Context,
			username string,
		) (*domain.User, error) {
			return &domain.User{
				ID:           1,
				Username:     "test",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	userService := NewUserServ(repo)

	req := domain.AuthRequest{
		Username: "test",
		Password: "test",
	}

	token, err := userService.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Error("expected token, got empty")
	}
}
func TestUserServ_Login_WrongPassword(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("another-password"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("can't hash password: %v", err)
	}
	repo := &MockUserRepo{
		FindByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
			if username != "test" {
				t.Errorf("expected username %q, got %q", "test", username)
			}
			return &domain.User{
				ID:           1,
				Username:     "test",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}
	userService := NewUserServ(repo)
	req := domain.AuthRequest{
		Username: "test",
		Password: "test",
	}
	token, err := userService.Login(context.Background(), req)
	if err == nil {
		t.Error("expected error")
	}
	if token != "" {
		t.Error("expected empty token")
	}
}
func TestUserServ_Login_UserNotFound(t *testing.T) {
	repo := &MockUserRepo{
		FindByUsernameFunc: func(
			ctx context.Context,
			username string,
		) (*domain.User, error) {
			return nil, sql.ErrNoRows
		},
	}

	userService := NewUserServ(repo)

	req := domain.AuthRequest{
		Username: "unknown",
		Password: "test",
	}

	token, err := userService.Login(context.Background(), req)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}

	if token != "" {
		t.Errorf("expected empty token, got %q", token)
	}
}
