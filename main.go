package main

import (
	"context"
	"errors"
	"log"
	"os"
	"task_manager/config"
	"task_manager/controllers"
	"task_manager/data"
	"task_manager/router"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	client, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal(errors.New("DB_NAME environment variable not set"))
	}
	db := client.Database(dbName)
	log.Printf("Connected to MongoDB database: %s", dbName)

	// TASK SERVICE (MongoDB)
	taskService := data.NewTaskService(db.Collection("tasks"))

	// USER SERVICE (MongoDB version needed)
	userService := data.NewUserService(db.Collection("users")) // <-- Must be Mongo-compatible

	taskController := controllers.NewTaskController(userService, taskService)
	r := router.SetupRouter(taskController)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
