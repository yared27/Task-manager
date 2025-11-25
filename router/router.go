package router

import (
	"os"
	"task_manager/config"
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	godotenv.Load()
	client, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	db := client.Database(os.Getenv("DB_NAME"))
	taskCollection := db.Collection("Tasks")
	taskService := data.NewTaskService(taskCollection)
	taskController := controllers.NewTaskController(taskService)

	api := r.Group("/tasks")
	{
		api.GET("/", taskController.GetAllTasks)
		api.GET("/:id", taskController.GetTaskByID)
		api.POST("/", taskController.CreateTask)
		api.PUT("/:id", taskController.UpdateTask)
		api.DELETE("/:id", taskController.DeleteTask)
	}

	return r
}
