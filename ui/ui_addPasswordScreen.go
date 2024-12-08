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

func showAddPasswordScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window) {
		var content *widget.Card
		var body *fyne.Container
		var viewOptiones = true

		returnButton := widget.NewButton(localizer.Get("return"), func() {
			controller.ShowScreen("main")
		})

		labelEntry := widget.NewEntry()
		labelEntry.SetPlaceHolder(localizer.Get("label"))

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder(localizer.Get("password"))

		lengthEntry := widget.NewEntry()
		lengthEntry.SetPlaceHolder(localizer.Get("lenght"))
		lengthEntry.Hide()

		useUpperCheck := widget.NewCheck(localizer.Get("includeUppercase"), nil)
		useLowerCheck := widget.NewCheck(localizer.Get("includeLowercase"), nil)
		useNumbersCheck := widget.NewCheck(localizer.Get("includeNumber"), nil)
		useSpecialsCheck := widget.NewCheck(localizer.Get("includeSpecialCharacter"), nil)
		useUpperCheck.Hide()
		useLowerCheck.Hide()
		useNumbersCheck.Hide()
		useSpecialsCheck.Hide()

		generateButton := widget.NewButton(localizer.Get("generateSafePassword"), func() {
			length := 16
			if lengthText := lengthEntry.Text; lengthText != "" {
				fmt.Sscanf(lengthText, "%d", &length)
			}

			password, err := contUser.GenerateNewPasswordSafe(length,
				useUpperCheck.Checked,
				useLowerCheck.Checked,
				useNumbersCheck.Checked,
				useSpecialsCheck.Checked)

			if err != nil {
				fmt.Println("Error generando la contraseña:", err)
			} else {
				passwordEntry.SetText(password)
			}
		})
		generateButton.Hide()

		labelErr := canvas.NewText(localizer.Get("addFailed"), theme.ErrorColor())
		labelErr.Hide()
		addButton := widget.NewButton(localizer.Get("add"), func() {
			_, err := dbController.InsertPassword(contUser.GetCurrentUserId(), labelEntry.Text, passwordEntry.Text, contUser.GetCurrentUserPassword())
			if err != nil {
				labelErr.Show()
				log.Println("Err in add: ", err)
				go func() {
					time.Sleep(5 * time.Second)
					labelErr.Hide()
				}()
			} else {
				controller.ShowScreen("main")
			}
		})

		var OptionsButton *widget.Button
		OptionsButton = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
			if viewOptiones {
				lengthEntry.Show()
				useUpperCheck.Show()
				useLowerCheck.Show()
				useNumbersCheck.Show()
				useSpecialsCheck.Show()
				generateButton.Show()
				OptionsButton.SetIcon(theme.MenuDropUpIcon())
			} else {
				lengthEntry.Hide()
				useUpperCheck.Hide()
				useLowerCheck.Hide()
				useNumbersCheck.Hide()
				useSpecialsCheck.Hide()
				generateButton.Hide()
				OptionsButton.SetIcon(theme.MenuDropDownIcon())
			}
			viewOptiones = !viewOptiones
			body.Refresh()
		})

		body = container.NewVBox(
			labelEntry,
			passwordEntry,
			lengthEntry,
			useUpperCheck,
			useLowerCheck,
			useNumbersCheck,
			useSpecialsCheck,
			generateButton,
			OptionsButton,
			addButton,
			returnButton,
			labelErr,
		)
		content = widget.NewCard(localizer.Get("addPassword"), localizer.Get("enterLabel&Password"), body)

		w.SetContent(content)
	}
}
