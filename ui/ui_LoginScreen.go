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
	"password_manager/localization"
	"time"
)

func showLoginScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window) {
		mailEntry := widget.NewEntry()
		mailEntry.SetPlaceHolder(localizer.Get("mail"))

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder(localizer.Get("password"))

		labelErr := canvas.NewText(localizer.Get("loginFailed"), theme.ErrorColor())
		labelErr.Hide()

		loginFunc := func() {
			err := contUser.Login(mailEntry.Text, passwordEntry.Text)
			if err != nil {
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

		loginButton := widget.NewButton(localizer.Get("login"), loginFunc)

		registerButton := widget.NewButton(localizer.Get("register"), func() {
			controller.ShowScreen("register")
		})

		form := widget.NewCard(localizer.Get("login"), localizer.Get("enterMail&Password"), container.NewVBox(
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
