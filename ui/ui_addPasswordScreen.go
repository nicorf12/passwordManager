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
	return func(w fyne.Window, params ...interface{}) {
		var content *fyne.Container
		var viewOptiones = true
		var scrollContent *container.Scroll

		returnButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			controller.ShowScreen("main")
		})

		addTitle := widget.NewLabel(localizer.Get("addPassword"))
		addSubtitle := widget.NewLabel(localizer.Get("enterLabel&Password"))
		passwordLabel := widget.NewLabel(localizer.Get("label"))
		passwordName := widget.NewLabel(localizer.Get("name"))
		passwordPassword := widget.NewLabel(localizer.Get("password"))
		passwordWebsite := widget.NewLabel(localizer.Get("website"))
		passwordFolder := widget.NewLabel(localizer.Get("folder"))
		passwordEncryption := widget.NewLabel(localizer.Get("encryption"))
		passwordNote := widget.NewLabel(localizer.Get("note"))

		labelEntry := widget.NewEntry()
		nameEntry := widget.NewEntry()
		websiteEntry := widget.NewEntry()
		passwordEntry := widget.NewPasswordEntry()

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

		folders, err := dbController.GetAllFolders()
		if err != nil {
			log.Fatal(err)
		}
		var folderKeys []string
		for folderName := range folders {
			folderKeys = append(folderKeys, folderName)
		}
		selectFolder := widget.NewSelect(folderKeys, func(selected string) { fmt.Println("Selected:", selected) })

		encryption, err := dbController.GetAllEncrypted()
		if err != nil {
			log.Fatal(err)
		}
		var encryptionKeys []string
		for encryptionName := range encryption {
			encryptionKeys = append(encryptionKeys, encryptionName)
		}
		selectEncryption := widget.NewSelect(encryptionKeys, func(selected string) { fmt.Println("Selected:", selected) })
		if len(encryptionKeys) > 0 {
			selectEncryption.Selected = encryptionKeys[0]
		}

		note := widget.NewMultiLineEntry()

		labelErr := canvas.NewText(localizer.Get("addFailed"), theme.ErrorColor())
		labelErr.Hide()
		addButton := widget.NewButton(localizer.Get("add"), func() {
			var idFolder int64
			if _, ok := folders[selectFolder.Selected]; ok {
				idFolder = folders[selectFolder.Selected]
			}
			var idEncryption int64
			if _, ok := encryption[selectEncryption.Selected]; ok {
				idEncryption = encryption[selectEncryption.Selected]
			}
			_, err := dbController.InsertPassword(
				contUser.GetCurrentUserId(),
				idFolder,
				labelEntry.Text,
				nameEntry.Text,
				passwordEntry.Text,
				websiteEntry.Text,
				note.Text,
				idEncryption,
				contUser.GetCurrentUserPassword())
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

		securityLevel := widget.NewProgressBar()
		securityLevel.Min = 0
		securityLevel.Max = 100
		evaluateButton := widget.NewButton(localizer.Get("evaluate"), func() {
			securityLevel.SetValue(contUser.GetPasswordSecurityLevel(passwordEntry.Text))
		})
		evaluateButton.Hide()
		securityLevel.Hide()

		var OptionsButton *widget.Button
		OptionsButton = widget.NewButtonWithIcon("", theme.MenuDropDownIcon(), func() {
			if viewOptiones {
				lengthEntry.Show()
				useUpperCheck.Show()
				useLowerCheck.Show()
				useNumbersCheck.Show()
				useSpecialsCheck.Show()
				generateButton.Show()
				evaluateButton.Show()
				securityLevel.Show()
				OptionsButton.SetIcon(theme.MenuDropUpIcon())
			} else {
				lengthEntry.Hide()
				useUpperCheck.Hide()
				useLowerCheck.Hide()
				useNumbersCheck.Hide()
				useSpecialsCheck.Hide()
				generateButton.Hide()
				evaluateButton.Hide()
				securityLevel.Hide()
				OptionsButton.SetIcon(theme.MenuDropDownIcon())
			}
			viewOptiones = !viewOptiones
			scrollContent.Refresh()
		})

		content = container.NewVBox(
			container.NewHBox(addTitle, returnButton),
			addSubtitle,
			passwordLabel,
			labelEntry,
			passwordName,
			nameEntry,
			passwordPassword,
			passwordEntry,
			lengthEntry,
			useUpperCheck,
			useLowerCheck,
			useNumbersCheck,
			useSpecialsCheck,
			generateButton,
			container.NewGridWithColumns(2, securityLevel, evaluateButton),
			OptionsButton,
			passwordWebsite,
			websiteEntry,
			passwordFolder,
			selectFolder,
			passwordEncryption,
			selectEncryption,
			passwordNote,
			note,
			addButton,
			labelErr,
		)

		scrollContent = container.NewScroll(content)
		w.SetContent(scrollContent)
	}
}
