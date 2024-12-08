package main

import (
	"log"
	"os"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/security"
	"password_manager/ui"
)

func main() {
	langEnv := os.Getenv("LANG")

	if langEnv == "" {
		langEnv = "en"
	}

	langCode := langEnv[:2]

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
		log.Println("Session found, starting authenticated UI.")
		contUser = controllers.NewControllerUserWithSession(dbController, session.UserID, session.UserMail, session.HashedPassword)
		ui.StartAuthenticatedUI(contUser, dbController, localizer)
	} else {
		log.Println("No session found, starting login UI.")
		contUser = controllers.NewControllerUser(dbController)
		ui.StartUI(contUser, dbController, localizer)
	}

	err = dbController.Close()
	if err != nil {
		log.Fatalf("Error closing connection to database: %v", err)
	}
}
