package main

import (
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/security"
	"password_manager/ui"
)

type AppContext struct {
	DBController *controllers.DBController
	Localizer    *localization.Localizer
	ContUser     *controllers.ControllerUser
}

func main() {
	appCtx, err := NewAppContext()
	if err != nil {
		log.Fatalf("Error inicializando la aplicación: %v", err)
	}
	ui.StartUI(appCtx.ContUser, appCtx.DBController, appCtx.Localizer)
}

func NewAppContext() (*AppContext, error) {
	dbController, err := controllers.NewDBController()
	if err != nil {
		return nil, err
	}

	config, _ := security.LoadConfig()
	if config == nil {
		config = security.LoadConfigDefault()
	}

	langCode := config.Lang
	localizer, err := localization.NewLocalizer(langCode)
	if err != nil {
		return nil, err
	}

	var contUser *controllers.ControllerUser
	session, err := security.LoadSession()
	if err == nil && session != nil {
		contUser, err = controllers.NewControllerUserWithSession(config, dbController, session.UserID, session.UserMail, session.HashedPassword)
		if err != nil {
			log.Println(err)
		}
	}

	if contUser == nil {
		contUser = controllers.NewControllerUser(config, dbController)
	}

	return &AppContext{dbController, localizer, contUser}, nil
}
