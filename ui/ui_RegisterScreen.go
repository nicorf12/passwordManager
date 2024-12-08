package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"time"
)

func showRegisterScreen(controller *controllers.ControllerScreen, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window) {
		var form *widget.Card
		emailEntry := widget.NewEntry()
		emailEntry.SetPlaceHolder(localizer.Get("mail"))

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder(localizer.Get("password"))

		labelErr := canvas.NewText("", theme.ErrorColor())
		labelErr.Hide()
		registerButton := widget.NewButton(localizer.Get("register"), func() {
			_, err := dbController.InsertUser(emailEntry.Text, passwordEntry.Text)
			if err != nil {
				labelErr.Text = localizer.Get("registerFailed")
				labelErr.Show()
				log.Println("Err in register: ", err)
				go func() {
					time.Sleep(5 * time.Second)
					labelErr.Hide()
				}()

				w.SetContent(form)
			} else {
				controller.ShowScreen("login")
			}
		})

		returnButton := widget.NewButton(localizer.Get("return"), func() {
			controller.ShowScreen("login")
		})

		form = widget.NewCard(localizer.Get("register"), localizer.Get("enterMail&PasswordToRegister"), container.NewVBox(
			emailEntry,
			passwordEntry,
			registerButton,
			returnButton,
			labelErr,
		))

		w.SetContent(form)
	}
}
