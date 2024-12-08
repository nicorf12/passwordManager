package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"time"
)

func showRegisterScreen(controller *controllers.ControllerScreen, dbController *controllers.DBController) controllers.Screen {
	return func(w fyne.Window) {
		var form *widget.Card
		emailEntry := widget.NewEntry()
		emailEntry.SetPlaceHolder("Email")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		labelErr := canvas.NewText("", theme.ErrorColor())
		labelErr.Hide()
		registerButton := widget.NewButton("Register", func() {
			_, err := dbController.InsertUser(emailEntry.Text, passwordEntry.Text)
			if err != nil {
				labelErr.Text = "Register failed: Incorrect mail or password"
				labelErr.Show()
				log.Println("Err in register: ", err)
				go func() {
					time.Sleep(5 * time.Second)
					labelErr.Hide()
				}()
			}
			w.SetContent(form)
		})

		returnButton := widget.NewButton("Return", func() {
			controller.ShowScreen("login")
		})

		form = widget.NewCard("Register", "Enter mail and password to register", container.NewVBox(
			emailEntry,
			passwordEntry,
			registerButton,
			returnButton,
			labelErr,
		))

		w.SetContent(form)
	}
}
