package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"time"
)

func showLoginScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser) controllers.Screen {
	return func(w fyne.Window) {
		mailEntry := widget.NewEntry()
		mailEntry.SetPlaceHolder("Mail")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		labelErr := canvas.NewText("", theme.ErrorColor())
		labelErr.Hide()

		loginFunc := func() {
			err := contUser.Login(mailEntry.Text, passwordEntry.Text)
			if err != nil {
				labelErr.Text = "Login failed: Incorrect mail or password"
				labelErr.Show()
				log.Println("Err in login: ", err)
				go func() {
					time.Sleep(5 * time.Second)
					labelErr.Hide()
				}()
			} else {
				controller.ShowScreen("main")
				fmt.Println("Logged successfully")
			}
		}

		loginButton := widget.NewButton("Login", loginFunc)

		registerButton := widget.NewButton("Register", func() {
			controller.ShowScreen("register")
		})

		form := widget.NewCard("Login", "Enter your mail and password", container.NewVBox(
			mailEntry,
			passwordEntry,
			loginButton,
			registerButton,
			labelErr))

		mailEntry.OnSubmitted = func(_ string) { loginFunc() }
		passwordEntry.OnSubmitted = func(_ string) { loginFunc() }

		w.SetContent(form)
	}
}
