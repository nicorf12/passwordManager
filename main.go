package main

import (
	"fmt"
	"fyne.io/systray"
	"log"
	"os"
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

	langCode := getLenguage()
	localizer, err := localization.NewLocalizer(langCode)
	if err != nil {
		return nil, err
	}

	var contUser *controllers.ControllerUser
	session, err := security.LoadSession()
	if err == nil && session != nil {
		fmt.Println("Session found, starting authenticated UI.")
		contUser = controllers.NewControllerUserWithSession(dbController, session.UserID, session.UserMail, session.HashedPassword)
	} else {
		fmt.Println("No session found, starting login UI.")
		contUser = controllers.NewControllerUser(dbController)
	}

	return &AppContext{dbController, localizer, contUser}, nil
}

func onReady(appCtx *AppContext, openUIChan chan bool) {
	systray.SetIcon(loadIcon("resources/dragon.ico"))
	systray.SetTitle("Gestor de Contraseñas")
	systray.SetTooltip("Gestor de Contraseñas en Ejecución")

	mOpen := systray.AddMenuItem("Abrir", "Abrir la aplicación")
	mQuit := systray.AddMenuItem("Salir", "Cerrar la aplicación")

	log.Println("Menú de bandeja agregado, esperando interacciones...")
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				log.Println("Se seleccionó 'Abrir'")
				openUIChan <- true
			case <-mQuit.ClickedCh:
				log.Println("Cerrando la aplicación...")
				systray.Quit()
			}
		}
	}()
}

func onExit(appCtx *AppContext) {
	log.Println("Aplicación cerrada.")
	if err := appCtx.DBController.Close(); err != nil {
		log.Fatalf("Error al cerrar la base de datos: %v", err)
	}
}

func getLenguage() string {
	return "es"
}

func loadIcon(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error al cargar el ícono: %v", err)
	}
	return data
}
