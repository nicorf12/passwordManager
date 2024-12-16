package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
)

func showFolderScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		var name string
		var folderId int64
		if len(params) > 1 {
			name, _ = params[0].(string)
			folderId, _ = params[1].(int64)
		}

		nameLabel := widget.NewLabel(name)
		deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			err := dbController.DeleteFolder(folderId)
			if err != nil {
				return
			}
			controller.ShowScreen("main")
		})

		passwords := container.NewVBox()
		var updatePasswordsList func()
		updatePasswordsList = func() {
			passwordsData, err := dbController.GetPasswordsByFolderAndUserID(contUser.GetCurrentUserId(), folderId, contUser.GetCurrentUserPassword())
			if err != nil {
				log.Printf("Error getting passwords: %v", err)
			} else {
				passwords.Objects = nil
				for _, password := range passwordsData {
					passwords.Add(createPasswordItem(w, controller, dbController, password, "folder", &name, &folderId, true))
				}
			}
		}

		updatePasswordsList()

		content := container.NewHSplit(menu(controller, contUser, dbController, localizer), container.NewVBox(
			container.NewHBox(nameLabel, deleteButton),
			passwords,
		))
		content.SetOffset(0.25)
		w.SetContent(content)
	}
}
