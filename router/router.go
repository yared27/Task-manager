package router
import (
	"github.com/gin-gonic/gin"
	"task_manager/controllers"
)

func SetupRouter() *gin.Engine{
	r := gin.Default()
	api := r.Group("/tasks")
	{
		api.GET("/", controllers.GetAllTasks)
		api.GET("/:id", controllers.GetTaskByID)
		api.POST("/",controllers.CreateTask)
		api.PUT("/:id", controllers.UpdateTask)
		api.DELETE("/:id", controllers.DeleteTask)
	}

	return  r
}