package handler

import (
	"Tamrin/tasks/internal/domain"
	"Tamrin/tasks/internal/service"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service service.IntService
}

func NewTaskHandler(sev service.IntService) *TaskHandler {
	return &TaskHandler{service: sev}
}
func (th *TaskHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	tasks, err := th.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}
func (th *TaskHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	task, errGetBy := th.service.GetById(ctx, id)
	if errGetBy != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errGetBy.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
func (th *TaskHandler) Create(c *gin.Context) {
	var task domain.Task
	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	_, errCreate := th.service.Create(ctx, task)
	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errCreate.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": task})
}
func (th *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var task domain.Task
	err = c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	errUpdating := th.service.Update(ctx, task, id)
	if errUpdating != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errUpdating.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}
func (th *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	errDeleting := th.service.Delete(ctx, id)
	if errors.Is(errDeleting, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
	}
	if errDeleting != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDeleting.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
