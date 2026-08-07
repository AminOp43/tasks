package service

import (
	"Tamrin/tasks/internal/domain"
	"Tamrin/tasks/internal/repository"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserServ struct {
	userRepo repository.UserRepository
}

func NewUserServ(userRepo repository.UserRepository) *UserServ {
	return &UserServ{userRepo}
}
func (u *UserServ) SignUp(ctx context.Context, req domain.AuthRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := domain.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}
	_, err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserServ) Login(ctx context.Context, req domain.AuthRequest) (string, error) {
	user, err := u.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	newJwt := jwt.MapClaims{"user_id": user.ID, "exp": expirationTime.Unix()}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newJwt)
	secret := []byte("my-jwt-secret")
	jwtTokenString, errJwt := jwtToken.SignedString(secret)
	if errJwt != nil {
		return "", errJwt
	}
	return jwtTokenString, nil
}
