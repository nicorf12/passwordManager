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
	"time"
)

func showAddPasswordScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController) controllers.Screen {
	return func(w fyne.Window) {
		var content *widget.Card
		var body *fyne.Container
		var viewOptiones = true

		returnButton := widget.NewButton("Return", func() {
			controller.ShowScreen("main")
		})

		labelEntry := widget.NewEntry()
		labelEntry.SetPlaceHolder("Label")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		lengthEntry := widget.NewEntry()
		lengthEntry.SetPlaceHolder("Length (e.g., 16)")
		lengthEntry.Hide()

		useUpperCheck := widget.NewCheck("Include Uppercase", nil)
		useLowerCheck := widget.NewCheck("Include Lowercase", nil)
		useNumbersCheck := widget.NewCheck("Include Numbers", nil)
		useSpecialsCheck := widget.NewCheck("Include Special Characters", nil)
		useUpperCheck.Hide()
		useLowerCheck.Hide()
		useNumbersCheck.Hide()
		useSpecialsCheck.Hide()

		generateButton := widget.NewButton("Generate Secure Password", func() {
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

		labelErr := canvas.NewText("", theme.ErrorColor())
		labelErr.Hide()
		addButton := widget.NewButton("Add", func() {
			_, err := dbController.InsertPassword(contUser.GetCurrentUserId(), labelEntry.Text, passwordEntry.Text, contUser.GetCurrentUserPassword())
			if err != nil {
				labelErr.Text = "Add failed: Incorrect label or password"
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
		content = widget.NewCard("Add Password", "Enter a label and password you want to register", body)

		w.SetContent(content)
	}
}
