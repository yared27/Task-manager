package router

import (
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(ctrl *controllers.TaskController) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", ctrl.Register)
	r.POST("/login", ctrl.Login)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/tasks", ctrl.GetAllTasks)
	auth.GET("/tasks/:id", ctrl.GetTaskByID)

	// Admin-only routes
	admin := auth.Group("/")
	admin.Use(middleware.AdminOnly())

	admin.POST("/tasks", ctrl.CreateTask)
	admin.PUT("/tasks/:id", ctrl.UpdateTask)
	admin.DELETE("/tasks/:id", ctrl.DeleteTask)
	admin.PUT("/promote/:id", ctrl.PromoteUser)

	return r
}
