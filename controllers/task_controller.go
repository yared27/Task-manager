package controllers

import (
	"net/http"
	"strconv"
	"task_manager/data"
	"task_manager/models"
	"time"
	"github.com/gin-gonic/gin"
)

// GetAllTasks handles GET /tasks
func GetAllTasks(c *gin.Context){
	tasks := data.ListTasks()
	c.JSON(http.StatusOK, gin.H{"tasks" : tasks})

}

// GetTaskByID handles GET /tasks/:id

func GetTaskByID(c *gin.Context){
	idStr := c.Param("id")
	id,err := strconv.Atoi(idStr)

	t, err := data.GetTask(id)
	if err != nil {
		if err == data.ErrorNotFound{
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK,t)
}

func CreateTask(c *gin.Context){
	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid input:" + err.Error()})
		return
	}

	due, err := time.Parse(time.RFC3339,input.DueDate)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "invalid due_date format; use RFC3339 (e.g. 2006-01-02T15:04:05z)"})
		return
	}
	t := data.CreateTask(input.Title, input.Description, due,input.Status)
	c.JSON(http.StatusCreated, t)

}

func UpdateTask(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid input:" + err.Error()})
		return
	}
	due, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "inavlid due_date format; use RFC3339 (e.g. 2006-01-02T15:04:05Z)"})
		return
	}
	updated,err := data.UpdateTask(id, input.Title, input.Description,due,input.Status)

	if err != nil{
		if err == data.ErrorNotFound{
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func DeleteTask(c *gin.Context){
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := data.DeleteTask(id)
	if err != nil{
		if err == data.ErrorNotFound{
		c.JSON(http.StatusNotFound, gin.H{"error": "task not foud"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error":"internal error"})
	return
	}
	c.Status(http.StatusNoContent)

}