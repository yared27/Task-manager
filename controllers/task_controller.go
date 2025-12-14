package controllers

import (
	"net/http"
	"time"

	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	UserService *data.UserService
	TaskService *data.TaskService
}

// ---------------- Constructor ----------------
func NewTaskController(userService *data.UserService, taskService *data.TaskService) *TaskController {
	return &TaskController{
		UserService: userService,
		TaskService: taskService,
	}
}

// ---------------- AUTH ----------------

// Register a new user
func (ctl *TaskController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.UserService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login user and return JWT token
func (ctl *TaskController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctl.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Promote user to admin (MongoDB)
func (ctl *TaskController) PromoteUser(c *gin.Context) {
	idHex := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := ctl.UserService.PromoteUser(userID); err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

// ---------------- TASKS ----------------

// Get all tasks
func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.TaskService.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// Get task by ID
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.TaskService.GetTask(id)
	if err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// Create a new task
func (tc *TaskController) CreateTask(c *gin.Context) {
	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	due, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format"})
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

	id, err := tc.TaskService.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

// Update a task
func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var input models.TaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	due, err := time.Parse(time.RFC3339, input.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format"})
		return
	}

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     due,
		Status:      input.Status,
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	if err := tc.TaskService.UpdateTask(id, &task); err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete a task
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := tc.TaskService.DeleteTask(id); err != nil {
		if err == data.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
