package main

import (
	"Tamrin/tasks/internal/handler"
	"Tamrin/tasks/internal/repository/postgres"
	"Tamrin/tasks/internal/service"
	"Tamrin/tasks/pkg/db"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	newDb := db.InitDB()
	repoNew := postgres.NewRepo(newDb)
	newTaskService := service.NewTaskService(repoNew)
	taskHandler := handler.NewTaskHandler(newTaskService)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	r.GET("/tasks", taskHandler.GetAll)
	r.POST("/tasks", taskHandler.Create)
	r.GET("/tasks/:id", taskHandler.GetById)
	r.PUT("/tasks/:id", taskHandler.Update)
	r.DELETE("/tasks/:id", taskHandler.Delete)
	r.Run(":8080")
}
