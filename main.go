package main

import (
	"log"
	"password_manager/internal/controllers"
	"password_manager/ui"
)

func main() {
	dbController, err := controllers.NewDBController()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	contUser := controllers.NewControllerUser(dbController)
	ui.StartUI(contUser, dbController)

	err = dbController.Close()
	if err != nil {
		log.Fatalf("Error closing connection to database: %v", err)
	}
}
