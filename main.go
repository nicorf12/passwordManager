package main

import (
	"log"
	"password_manager/internal/controllers"
	"password_manager/security"
	"password_manager/ui"
)

func main() {
	dbController, err := controllers.NewDBController()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	var contUser *controllers.ControllerUser
	session, err := security.LoadSession()
	if err == nil && session != nil {
		log.Println("Session found, starting authenticated UI.")
		contUser = controllers.NewControllerUserWithSession(dbController, session.UserID, session.UserMail, session.HashedPassword)
		ui.StartAuthenticatedUI(contUser, dbController)
	} else {
		log.Println("No session found, starting login UI.")
		contUser = controllers.NewControllerUser(dbController)
		ui.StartUI(contUser, dbController)
	}

	err = dbController.Close()
	if err != nil {
		log.Fatalf("Error closing connection to database: %v", err)
	}
}
