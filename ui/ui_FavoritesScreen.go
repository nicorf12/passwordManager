package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
)

func showFavoritesScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		favoritesTitle := widget.NewLabel(localizer.Get("favorites"))
		favoritesSubtitle := widget.NewLabel(localizer.Get("hereCanSeeFavoritesPasswords"))

		passwords := container.NewVBox()
		var updatePasswordsList func()
		updatePasswordsList = func() {
			passwordsData, err := dbController.GetPasswordsByFavoriteAndUserID(contUser.GetCurrentUserId(), contUser.GetCurrentUserPassword())
			if err != nil {
				log.Printf("Error getting passwords: %v", err)
			} else {
				passwords.Objects = nil
				for _, password := range passwordsData {

					passwords.Add(createPasswordItem(w, controller, dbController, password, "favorites", nil, nil, true))
				}
			}
		}

		updatePasswordsList()

		content := container.NewHSplit(menu(controller, contUser, dbController, localizer), container.NewVBox(
			favoritesTitle,
			favoritesSubtitle,
			passwords,
		))
		content.SetOffset(0.25)
		w.SetContent(content)
	}
}
