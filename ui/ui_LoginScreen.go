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
	"io"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"time"
)

var questions = []string{
	"¿Hola?",
	"¿Como estas?",
	"¿Bien?",
}

func showLoginScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, localizer *localization.Localizer, app fyne.App) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		mailEntry := widget.NewEntry()
		mailEntry.SetPlaceHolder(localizer.Get("mail"))

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder(localizer.Get("password"))

		labelErr := canvas.NewText(localizer.Get("loginFailed"), theme.Color(theme.ColorNameError))
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

		importErrorLabel := canvas.NewText(localizer.Get("importFailed"), theme.Color(theme.ColorNameError))
		importErrorLabel.Hide()

		importButton := widget.NewButton(localizer.Get("import"), func() {
			var wImport fyne.Window
			var data string
			wImport = windowImportExport(app, func(answer0, answer1, answer2 string) {
				startDir, err := storage.ListerForURI(storage.NewFileURI("."))
				if err != nil {
					log.Println("Error obteniendo directorio inicial:", err)
					importErrorLabel.Show()
					go func() {
						time.Sleep(5 * time.Second)
						importErrorLabel.Hide()
					}()
					return
				}
				fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if reader != nil {
						fmt.Println("Archivo seleccionado:", reader.URI().Path())
						defer reader.Close()
						fileData, readErr := io.ReadAll(reader)
						if readErr != nil {
							log.Println("Error leyendo el archivo:", readErr)
							importErrorLabel.Show()
							go func() {
								time.Sleep(5 * time.Second)
								importErrorLabel.Hide()
							}()
							return
						}
						data = string(fileData)
						data, err = contUser.DecryptToImport(data, answer0, answer1, answer2)
						if err != nil {
							log.Println("Error al importar:", err)
							importErrorLabel.Show()
							go func() {
								time.Sleep(5 * time.Second)
								importErrorLabel.Hide()
							}()
						} else {
							fmt.Println("Imported successfully")
							wImport.Close()
						}
					} else {
						log.Println("No se seleccionó ningún archivo")
						importErrorLabel.Show()
						go func() {
							time.Sleep(5 * time.Second)
							importErrorLabel.Hide()
						}()
					}
				}, wImport)
				fileDialog.SetLocation(startDir)
				fileDialog.Show()
			},
				localizer,
				"import",
				importErrorLabel)
			wImport.Show()
		})

		form := widget.NewCard(localizer.Get("login"), localizer.Get("enterMail&Password"), container.NewVBox(
			mailEntry,
			passwordEntry,
			loginButton,
			registerButton,
			labelErr,
			importButton))

		mailEntry.OnSubmitted = func(_ string) { loginFunc() }
		passwordEntry.OnSubmitted = func(_ string) { loginFunc() }

		w.SetContent(form)
	}
}
