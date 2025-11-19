package main

import (
	"log"
	"task_manager/router"

)



// var Tasks = []Task{
// 	{ID: "1", Title: "SOFTWARE ENGINEERING", Description: "full time full-stack", DueDate: time.Date(2025, time.December, 4, 0, 0, 0, 0, time.UTC), Status: "Done"},
// 	{ID: "2", Title: "SOFTWARE maintaniner", Description: "full time maintainer", DueDate: time.Date(2025, time.December, 4, 0, 0, 0, 0, time.UTC), Status: "InProgress"},
// 	{ID: "3", Title: "SOFTWARE Tester", Description: "full time software Tester", DueDate: time.Date(2025, time.December, 4, 0, 0, 0, 0, time.UTC), Status: "Done"},
// }

// func getAllTasks(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, Tasks)
// }

// func getTaskByID(c *gin.Context) {
// 	id := c.Param("id")
// 	for _, task := range Tasks {
// 		if task.ID == id {
// 			c.IndentedJSON(http.StatusOK, task)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
// }

// func updateTask(c *gin.Context) {
// 	id := c.Param("id")

// 	var updatedTask Task
// 	if err := c.ShouldBindJSON(&updatedTask); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	for i, task := range Tasks {
// 		if task.ID == id {
// 			Tasks[i].Title = updatedTask.Title
// 			Tasks[i].Description = updatedTask.Description
// 			Tasks[i].DueDate = updatedTask.DueDate
// 			Tasks[i].Status = updatedTask.Status
// 			c.IndentedJSON(http.StatusOK, Tasks[i])
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
// }
// func deleteTask(c *gin.Context){
// 	id := c.Param("id")

// 	for i, task := range Tasks{
// 		if task.ID == id{
// 			Tasks = append(Tasks[:i], Tasks[i+1:])
// 		}
// 	}
// }

func main() {
	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	
}