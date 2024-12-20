package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/ui/themes"
	_ "password_manager/ui/themes"
)

func StartUI(contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) {
	mainApp := app.NewWithID("com.example.passwordmanager")
	_, t := contUser.GetConfig()
	setSettings(mainApp, themes.GetTheme(t))
	mainWin := mainApp.NewWindow((*localizer).Get("passwordManager"))
	icono, err := fyne.LoadResourceFromPath("resources/dragon.png")
	if err != nil {
		log.Printf("Error cargando el ícono: %v", err)
		icono = theme.ErrorIcon()
	}
	mainWin.SetIcon(icono)
	mainWin.Resize(fyne.NewSize(800, 600))

	/*
		if desk, ok := mainApp.(desktop.App); ok {
			m := fyne.NewMenu(localizer.Get("passwordManager"),
				fyne.NewMenuItem(localizer.Get("open"), func() {
					mainWin.Show()
				}),
				fyne.NewMenuItem(localizer.Get("quit"), func() {
					err := dbController.Close() // arreglar para q no este aca
					if err != nil {
						return
					}
					mainApp.Quit()
				}),
			)
			trayIcon, _ := fyne.LoadResourceFromPath("resources/dragon.ico")
			desk.SetSystemTrayIcon(trayIcon)
			desk.SetSystemTrayMenu(m)
		}

		mainWin.SetCloseIntercept(func() {
			mainWin.Hide()
		})
	*/

	screenController := controllers.NewControllerScreen(mainWin)
	screenController.RegisterScreen("login", showLoginScreen(screenController, contUser, localizer, mainApp))
	screenController.RegisterScreen("main", showMainScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("register", showRegisterScreen(screenController, dbController, localizer))
	screenController.RegisterScreen("add", showAddPasswordScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("view", showViewScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("folder", showFolderScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("favorites", showFavoritesScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("config", showConfigScreen(screenController, contUser, dbController, localizer, mainApp))

	if contUser.SomeoneLoggedIn() {
		screenController.ShowScreen("main")
	} else {
		screenController.ShowScreen("login")
	}

	mainWin.ShowAndRun()
}

func setSettings(app fyne.App, theme fyne.Theme) {
	app.Settings().SetTheme(theme)
}
