package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
)

func showLoginScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser) controllers.Screen {
	return func(w fyne.Window) {
		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("User")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		loginButton := widget.NewButton("Login", func() {
			err := contUser.Login(usernameEntry.Text, passwordEntry.Text)
			if err != nil {
				log.Println("Err in login: ", err)
			} else {
				controller.ShowScreen("main")
				fmt.Println("Logged successfully")
			}
		})

		registerButton := widget.NewButton("Register", func() {
			controller.ShowScreen("register")
		})

		form := container.NewVBox(
			widget.NewLabel("Login"),
			usernameEntry,
			passwordEntry,
			loginButton,
			registerButton,
		)

		w.SetContent(form)
	}
}
