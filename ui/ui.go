package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"password_manager/internal/controllers"
)

func StartUI(contUser *controllers.ControllerUser, dbController *controllers.DBController) {
	a := app.New()
	a.Settings().SetTheme(&customDarkTheme{})
	w := a.NewWindow("Password Manager")

	w.Resize(fyne.NewSize(400, 400))
	screenController := controllers.NewControllerScreen(w)

	screenController.RegisterScreen("login", showLoginScreen(screenController, contUser))
	screenController.RegisterScreen("main", showMainScreen(screenController, contUser, dbController))
	screenController.RegisterScreen("register", showRegisterScreen(screenController, dbController))
	screenController.RegisterScreen("add", showAddPasswordScreen(screenController, contUser, dbController))

	screenController.ShowScreen("login")
	w.ShowAndRun()
}

func StartAuthenticatedUI(contUser *controllers.ControllerUser, dbController *controllers.DBController) {
	a := app.New()
	a.Settings().SetTheme(&customDarkTheme{})
	w := a.NewWindow("Password Manager")

	w.Resize(fyne.NewSize(400, 400))
	screenController := controllers.NewControllerScreen(w)

	screenController.RegisterScreen("login", showLoginScreen(screenController, contUser))
	screenController.RegisterScreen("main", showMainScreen(screenController, contUser, dbController))
	screenController.RegisterScreen("register", showRegisterScreen(screenController, dbController))
	screenController.RegisterScreen("add", showAddPasswordScreen(screenController, contUser, dbController))

	screenController.ShowScreen("main")
	w.ShowAndRun()
}
