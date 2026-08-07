package main

import (
	""
	"Tamrin/tasks/internal/handler"
	"Tamrin/tasks/internal/middleware"
	"Tamrin/tasks/internal/repository/postgres"
	"Tamrin/tasks/internal/service"
	"Tamrin/tasks/pkg/db"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	newDb := db.InitDB()
	taskRepo := postgres.NewRepo(newDb)
	userRepo := postgres.NewUserRepo(newDb)

	newTaskService := service.NewTaskService(taskRepo)
	newUserService := service.NewUserServ(userRepo)

	taskHandler := handler.NewTaskHandler(newTaskService)
	userHandler := handler.NewUserHandler(newUserService)

	r := gin.Default()
	tasks := r.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware)

	tasks.GET("", taskHandler.GetAll)
	tasks.POST("", taskHandler.Create)
	tasks.GET("/:id", taskHandler.GetById)
	tasks.PUT("/:id", taskHandler.Update)
	tasks.DELETE("/:id", taskHandler.Delete)

	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	r.POST("/user/signup", userHandler.SignUp)
	r.POST("/user/login", userHandler.Login)

	r.Run(":8080")
}
