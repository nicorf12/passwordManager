package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"password_manager/internal/controllers"
	"password_manager/localization"
)

func StartUI(contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) {
	a := app.New()
	a.Settings().SetTheme(&customDarkTheme{})
	w := a.NewWindow(localizer.Get("passwordManager"))
	icono, err := fyne.LoadResourceFromPath("resources/icono_window.png")
	if err != nil {
		panic(err)
	}
	w.SetIcon(icono)

	w.Resize(fyne.NewSize(400, 400))
	screenController := controllers.NewControllerScreen(w)

	screenController.RegisterScreen("login", showLoginScreen(screenController, contUser, localizer))
	screenController.RegisterScreen("main", showMainScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("register", showRegisterScreen(screenController, dbController, localizer))
	screenController.RegisterScreen("add", showAddPasswordScreen(screenController, contUser, dbController, localizer))

	screenController.ShowScreen("login")
	w.ShowAndRun()
}

func StartAuthenticatedUI(contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) {
	a := app.New()
	a.Settings().SetTheme(&customDarkTheme{})
	w := a.NewWindow(localizer.Get("passwordManager"))
	icono, err := fyne.LoadResourceFromPath("resources/icono_window.png")
	if err != nil {
		panic(err)
	}
	w.SetIcon(icono)

	w.Resize(fyne.NewSize(400, 400))
	screenController := controllers.NewControllerScreen(w)

	screenController.RegisterScreen("login", showLoginScreen(screenController, contUser, localizer))
	screenController.RegisterScreen("main", showMainScreen(screenController, contUser, dbController, localizer))
	screenController.RegisterScreen("register", showRegisterScreen(screenController, dbController, localizer))
	screenController.RegisterScreen("add", showAddPasswordScreen(screenController, contUser, dbController, localizer))

	screenController.ShowScreen("main")
	w.ShowAndRun()
}
