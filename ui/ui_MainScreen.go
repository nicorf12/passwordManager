package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"strconv"
)

func showMainScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		passwordsTitle := widget.NewLabel(localizer.Get("passwords"))
		passwordsSubtitle := widget.NewLabel(localizer.Get("hereCanSeePasswords"))

		addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
			controller.ShowScreen("add")
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

					passwords.Add(createPasswordItem(w, controller, dbController, password, "main", nil, nil, true))
				}
			}
		}

		updatePasswordsList()

		content := container.NewHSplit(menu(controller, contUser, dbController, localizer), container.NewVBox(
			container.NewHBox(passwordsTitle, addButton),
			passwordsSubtitle,
			passwords,
		))
		content.SetOffset(0.18)
		w.SetContent(content)
	}
}

func createPasswordItem(w fyne.Window, controller *controllers.ControllerScreen, dbController *controllers.DBController, passwordDetails map[string]string, returnScreen string, folderName *string, folderId *int64, isDarkTheme bool) fyne.CanvasObject {
	icon := GetIconForPassword(passwordDetails, isDarkTheme)
	iconWidget := widget.NewIcon(icon)
	iconWidget.Resize(fyne.NewSize(60, 60))
	iconContainer := container.New(
		layout.NewGridLayout(3),
		layout.NewSpacer(),
		iconWidget,
		layout.NewSpacer(),
	)
	iconContainer.Resize(fyne.NewSize(80, 80))

	nameLabel := widget.NewLabel(passwordDetails["label"])
	detailLabel := widget.NewLabel(passwordDetails["name"])

	var iconFavorite fyne.Resource
	if passwordDetails["is_favorite"] == "1" {
		iconFavorite = theme.CheckButtonCheckedIcon()
	} else {
		iconFavorite = theme.CheckButtonIcon()
	}

	var favoriteButton *widget.Button
	favoriteButton = widget.NewButtonWithIcon("", iconFavorite, func() {
		passwordID, _ := strconv.ParseInt(passwordDetails["id"], 10, 64)
		dbController.EditFavoritePassword(passwordID)
		if passwordDetails["is_favorite"] == "1" {
			passwordDetails["is_favorite"] = "0"
			favoriteButton.SetIcon(theme.CheckButtonIcon())
		} else {
			passwordDetails["is_favorite"] = "1"
			favoriteButton.SetIcon(theme.CheckButtonCheckedIcon())
		}
		favoriteButton.Refresh()
	})
	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(passwordDetails["password"])
	})
	detailButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		if returnScreen == "main" || returnScreen == "favorites" {
			controller.ShowScreen("view", passwordDetails, returnScreen)
		} else {
			controller.ShowScreen("view", passwordDetails, returnScreen, *folderName, *folderId)
		}
	})

	info := container.NewVBox(
		nameLabel,
		detailLabel,
	)

	buttons := container.NewHBox(
		favoriteButton,
		copyButton,
		detailButton,
	)

	return container.NewBorder(nil, nil, iconContainer, buttons, info)
}
