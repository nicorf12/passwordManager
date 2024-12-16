package main

import (
	"fmt"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/security"
	"password_manager/ui"
)

func main() {
	langCode := getLenguage()

	dbController, err := controllers.NewDBController()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	localizer, err := localization.NewLocalizer(langCode)
	if err != nil {
		log.Fatalf("Error initializing localizer: %v", err)
	}

	var contUser *controllers.ControllerUser
	session, err := security.LoadSession()
	if err == nil && session != nil {
		fmt.Println("Session found, starting authenticated UI.")
		contUser = controllers.NewControllerUserWithSession(dbController, session.UserID, session.UserMail, session.HashedPassword)
		ui.StartAuthenticatedUI(contUser, dbController, localizer)
	} else {
		fmt.Println("No session found, starting login UI.")
		contUser = controllers.NewControllerUser(dbController)
		ui.StartUI(contUser, dbController, localizer)
	}

	err = dbController.Close()
	if err != nil {
		log.Fatalf("Error closing connection to database: %v", err)
	}
}

func getLenguage() string {
	return "es"
}
