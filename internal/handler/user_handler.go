package handler

import (
	"Tamrin/tasks/internal/domain"
	"Tamrin/tasks/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(userHandler service.UserService) *UserHandler {
	return &UserHandler{service: userHandler}
}
func (u *UserHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	var user domain.AuthRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = u.service.SignUp(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}
func (u *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var user domain.AuthRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, errLog := u.service.Login(ctx, user)
	if errLog != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errLog.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}
