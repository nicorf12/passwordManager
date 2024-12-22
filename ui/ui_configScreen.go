package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"password_manager/ui/themes"
	"time"
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
			setSettings(app, themes.GetTheme(themeSelected))
			controller.ShowScreen("config")
		})

		returnButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			controller.ShowScreen("main")
		})

		exportErrorLabel := canvas.NewText(localizer.Get("exportFailed"), theme.Color(theme.ColorNameError))
		exportErrorLabel.Hide()

		exportButton := widget.NewButton(localizer.Get("export"), func() {
			data, err := dbController.GetDataToExport(contUser.GetCurrentUserId())
			if err != nil {
				fmt.Println("error extracting data")
			}
			var wExport fyne.Window
			wExport = windowImportExport(
				app,
				func(answer0, answer1, answer2 string) {
					data, err := contUser.EncryptToExport(data, answer0, answer1, answer2)
					if err != nil {
						fmt.Println("Error encrypting data:", err)
					}
					startDir, err := storage.ListerForURI(storage.NewFileURI("."))
					if err != nil {
						log.Println("Error obteniendo directorio inicial:", err)
						return
					}
					fileDialog := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
						if err != nil {
							exportErrorLabel.Show()
							go func() {
								time.Sleep(5 * time.Second)
								exportErrorLabel.Hide()
							}()
							return
						}
						if uri == nil {
							fmt.Println("File not selected")
						} else {
							fmt.Fprintln(uri, data)
							uri.Close()
							wExport.Close()
						}
					}, wExport)
					fileDialog.SetLocation(startDir)
					fileDialog.Show()
				},
				localizer,
				"export",
				exportErrorLabel,
			)
			wExport.Show()
		})

		content := container.NewVBox(
			selectLang,
			selectTheme,
			saveButton,
			exportButton,
			returnButton,
		)
		w.SetContent(content)
	}
}
