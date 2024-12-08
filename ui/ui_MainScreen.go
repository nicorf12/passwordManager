package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"strconv"
)

func showMainScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window) {
		logoutButton := widget.NewButton(localizer.Get("logout"), func() {
			contUser.Logout()
			controller.ShowScreen("login")
		})

		passwords := container.NewVBox()
		var updatePasswordsList func()
		updatePasswordsList = func() {
			passwordsData, err := dbController.GetPasswordsByUserID(contUser.GetCurrentUserId(), contUser.GetCurrentUserPassword())
			if err != nil {
				log.Printf("Error getting passwords: %v", err)
			} else {
				passwords.Objects = nil
				for _, password := range passwordsData {
					var buttonBar *fyne.Container

					passwordEntry := widget.NewEntry()
					passwordEntry.Disable()
					passwordEntry.Hide()
					passwordEntry.SetText(password["password"])

					var showHideButton *widget.Button
					isVisible := false
					showHideButton = widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
						if isVisible {
							passwordEntry.Hide()
							showHideButton.SetIcon(theme.VisibilityIcon())
							passwordEntry.Refresh()
						} else {
							passwordEntry.Show()
							showHideButton.SetIcon(theme.VisibilityOffIcon())
							passwordEntry.Refresh()
						}
						isVisible = !isVisible
					})

					var deleteButton *widget.Button
					var dialog *widget.PopUp
					deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
						dialog = widget.NewPopUp(container.NewVBox(
							widget.NewLabel(localizer.Get("sureToDeletePassword")),
							container.NewHBox(
								widget.NewButton(localizer.Get("cancel"), func() {
									dialog.Hide()
								}),
								widget.NewButton("OK", func() {
									passwordID, _ := strconv.ParseInt(password["id"], 10, 64)
									err := dbController.DeletePassword(passwordID)
									if err != nil {
										fmt.Printf("Error deleting password: %v", err)
									} else {
										fmt.Printf("Deleted password: %v:%v", password["label"], password["password"])
									}
									updatePasswordsList()
									dialog.Hide()
								}),
							),
						), w.Canvas())
						dialog.Show()
					})

					var confirmEditionButton *widget.Button
					var cancelEditionButton *widget.Button
					var editButton *widget.Button
					editButton = widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
						confirmEditionButton.Show()
						cancelEditionButton.Show()
						passwordEntry.Show()
						if !isVisible {
							isVisible = true
							showHideButton.SetIcon(theme.VisibilityOffIcon())
						}
						editButton.Hide()
						passwordEntry.Enable()
					})

					confirmEditionButton = widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
						passwordID, _ := strconv.ParseInt(password["id"], 10, 64)
						err := dbController.EditPassword(passwordID, passwordEntry.Text, contUser.GetCurrentUserPassword())
						if err != nil {
							fmt.Println("Edited password: short password")
						} else {
							fmt.Println("Edited password: %v:%v", password["password"], passwordEntry.Text)
						}
						editButton.Show()
						confirmEditionButton.Hide()
						cancelEditionButton.Hide()
						passwordEntry.Disable()
						updatePasswordsList()
					})
					confirmEditionButton.Hide()

					cancelEditionButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
						passwordEntry.SetText(password["password"])
						editButton.Show()
						confirmEditionButton.Hide()
						cancelEditionButton.Hide()
						passwordEntry.Disable()
					})
					cancelEditionButton.Hide()

					copyButton := widget.NewButtonWithIcon(localizer.Get("copy"), theme.ContentCopyIcon(), func() {
						textToCopy := passwordEntry.Text
						w.Clipboard().SetContent(textToCopy)
					})

					labelEntry := widget.NewLabel(password["label"])
					buttonBar = container.NewVBox(container.NewHBox(labelEntry, showHideButton, deleteButton, editButton, confirmEditionButton, cancelEditionButton, copyButton), container.NewVBox(passwordEntry))
					passwords.Add(buttonBar)
				}
			}
		}

		updatePasswordsList()

		addButton := widget.NewButton(localizer.Get("add"), func() {
			controller.ShowScreen("add")
		})

		content := widget.NewCard(localizer.Get("passwords"), localizer.Get("hereCanSeePasswords"), container.NewVBox(
			passwords,
			addButton,
			logoutButton,
		))

		w.SetContent(content)
	}
}
