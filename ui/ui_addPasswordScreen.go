package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"password_manager/internal/controllers"
)

func showAddPasswordScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController) controllers.Screen {
	return func(w fyne.Window) {
		returnButton := widget.NewButton("Return", func() {
			controller.ShowScreen("main")
		})

		labelEntry := widget.NewEntry()
		labelEntry.SetPlaceHolder("Label")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		addButton := widget.NewButton("Add", func() {
			_, err := dbController.InsertPassword(contUser.GetCurrentUserId(), labelEntry.Text, passwordEntry.Text, contUser.GetCurrentUserPassword())
			if err != nil {
				fmt.Println(err)
			} else {
				controller.ShowScreen("main")
			}
		})

		content := container.NewVBox(
			widget.NewLabel("Add Password"),
			labelEntry,
			passwordEntry,
			addButton,
			returnButton,
		)

		w.SetContent(content)
	}
}
