package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/ui/themes"
)

var (
	languagesNames = []string{
		"en",
		"es",
		"de",
		"fr",
		"it",
		"pt",
	}
	themesNames = []string{
		"Dark",
		"Light",
		"Pink",
		"Blue",
	}
)

func showConfigScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer, app fyne.App) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		var selectLang *widget.Select
		var selectTheme *widget.Select
		var langSelected string
		var themeSelected string

		currentLang, currentTheme := contUser.GetConfig()

		selectLang = widget.NewSelect(languagesNames, func(s string) {
			langSelected = selectLang.Selected
		})
		selectLang.Selected = currentLang
		langSelected = currentLang
		selectTheme = widget.NewSelect(themesNames, func(s string) {
			themeSelected = selectTheme.Selected
		})
		selectTheme.Selected = currentTheme
		themeSelected = currentTheme

		saveButton := widget.NewButton(localizer.Get("save"), func() {
			contUser.SetConfig(langSelected, themeSelected)
			err := localizer.UpdateTranslations(langSelected)
			if err != nil {
				fmt.Println("Error updating translations:", err)
			}
			setSettings(app, GetTheme(themeSelected))
			controller.ShowScreen("config")
		})

		returnButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			controller.ShowScreen("main")
		})

		content := container.NewVBox(selectLang, selectTheme, saveButton, returnButton)
		w.SetContent(content)
	}
}

func GetTheme(theme string) fyne.Theme {
	switch theme {
	case "Dark":
		return themes.CustomDarkTheme{}
	case "Light":
		return themes.CustomLightTheme{}
	case "Pink":
		return themes.CustomPinkTheme{}
	case "Blue":
		return themes.CustomBlueTheme{}
	default:
		return themes.CustomDarkTheme{}
	}
}
