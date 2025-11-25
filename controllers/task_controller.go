package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"task_manager/data"
	"task_manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	Service *data.TaskService
}

// Constructor
func NewTaskController(s *data.TaskService) *TaskController {
	return &TaskController{Service: s}
}

// ------------------------------------------------------------
// GET /tasks
// ------------------------------------------------------------
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.Service.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// ------------------------------------------------------------
// GET /tasks/:id
// ------------------------------------------------------------
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.Service.GetTask(id)
	if err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// ------------------------------------------------------------
// POST /tasks
// ------------------------------------------------------------
func (tc *TaskController) CreateTask(c *gin.Context) {
	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Parse due date
	due, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format (use RFC3339)"})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     due,
		Status:      input.Status,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	result, err := tc.Service.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created",
		"id":      result.Hex(),
	})
}

// ------------------------------------------------------------
// PUT /tasks/:id
// ------------------------------------------------------------
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	due, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format (use RFC3339)"})
		return
	}
	updatedTask := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     due,
		Status:      input.Status,
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	err = tc.Service.UpdateTask(id, &updatedTask)
	if err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": updatedTask})
}

// ------------------------------------------------------------
// DELETE /tasks/:id
// ------------------------------------------------------------
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := tc.Service.DeleteTask(id)
	if err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}
