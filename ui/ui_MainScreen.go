package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
)

func showMainScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController) controllers.Screen {
	return func(w fyne.Window) {
		logoutButton := widget.NewButton("Logout", func() {
			contUser.Logout()
			controller.ShowScreen("login")
		})

		passwords := container.NewVBox()
		passwordsData, err := dbController.GetPasswordsByUserID(contUser.GetCurrentUserId(), contUser.GetCurrentUserPassword())
		if err != nil {
			log.Printf("Error getting passwords: %v", err)
		} else {
			for _, password := range passwordsData {
				passwordEntry := widget.NewEntry()
				passwordEntry.Disable()
				passwordEntry.Hide()
				passwordEntry.SetText(password["password"])

				var showHideButton *widget.Button
				showHideButton = widget.NewButton("Show", func() {
					if showHideButton.Text == "Hide" {
						passwordEntry.Hide()
						passwordEntry.Refresh()
						showHideButton.SetText("Show")
					} else {
						passwordEntry.Show()
						passwordEntry.Refresh()
						showHideButton.SetText("Hide")
					}
				})
				labelEntry := widget.NewLabel(password["label"])

				passwords.Add(container.NewVBox(container.NewHBox(labelEntry, showHideButton), container.NewVBox(passwordEntry)))
			}
		}

		addButton := widget.NewButton("Add", func() {
			controller.ShowScreen("add")
		})

		content := container.NewVBox(
			widget.NewLabel("Passwords"),
			widget.NewLabel("My mail: "+contUser.GetCurrentUserEmail()),
			passwords,
			addButton,
			logoutButton,
		)

		w.SetContent(content)
	}
}
