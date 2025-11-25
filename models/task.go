package models

import ("time"
		"go.mongodb.org/mongo-driver/bson/primitive")

type Task struct {
	ID          primitive.ObjectID    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	CreatedAt 	primitive.DateTime   `bson:"created_at"`
	UpdatedAt 	primitive.DateTime   `bson:"updated_at"`
}

type TaskInput struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
    DueDate     string `json:"due_date" binding:"required"`
    Status      string `json:"status" binding:"required"`
}
