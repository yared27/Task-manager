package main

import (
	"log"
	"task_manager/router"
	"task_manager/config"

)

func main() {
	config.ConnectDB()
	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	
}