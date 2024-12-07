package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"password_manager/internal/controllers"
)

func showRegisterScreen(controller *controllers.ControllerScreen, dbController *controllers.DBController) controllers.Screen {
	return func(w fyne.Window) {
		var form *fyne.Container
		emailEntry := widget.NewEntry()
		emailEntry.SetPlaceHolder("Email")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		estado := ""
		registerButton := widget.NewButton("Register", func() {
			_, err := dbController.InsertUser(emailEntry.Text, passwordEntry.Text)
			if err != nil {
				estado = "Email already registered"
			} else {
				estado = "Registration successful"
			}
			w.SetContent(form)
		})

		returnButton := widget.NewButton("Return", func() {
			controller.ShowScreen("login")
		})

		estadoRegistro := widget.NewLabel(estado)

		form = container.NewVBox(
			widget.NewLabel("Register"),
			emailEntry,
			passwordEntry,
			registerButton,
			returnButton,
			estadoRegistro,
		)

		w.SetContent(form)
	}
}
