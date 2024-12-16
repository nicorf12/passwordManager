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

func showViewScreen(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) controllers.Screen {
	return func(w fyne.Window, params ...interface{}) {
		var password map[string]string
		var screenReturn string
		var folderName string
		var folderId int64
		password = params[0].(map[string]string)
		if len(params) > 1 {
			screenReturn = params[1].(string)
			if screenReturn == "folder" {
				folderName = params[2].(string)
				folderId = params[3].(int64)
			}
		}

		viewTitle := widget.NewLabel(localizer.Get("view"))
		passwordLabel := widget.NewLabel(localizer.Get("label"))
		passwordName := widget.NewLabel(localizer.Get("name"))
		passwordPassword := widget.NewLabel(localizer.Get("password"))
		passwordWebsite := widget.NewLabel(localizer.Get("website"))
		passwordFolder := widget.NewLabel(localizer.Get("folder"))
		passwordEncryption := widget.NewLabel(localizer.Get("encryption"))
		passwordNote := widget.NewLabel(localizer.Get("note"))

		labelEntry := widget.NewEntry()
		labelEntry.Disable()
		labelEntry.SetText(password["label"])

		nameEntry := widget.NewEntry()
		nameEntry.Disable()
		nameEntry.SetText(password["name"])

		websiteEntry := widget.NewEntry()
		websiteEntry.Disable()
		websiteEntry.SetText(password["website"])

		passwordEntry := widget.NewEntry()
		passwordEntry.Disable()
		passwordEntry.SetText("")

		var showHideButton *widget.Button
		isVisible := false
		showHideButton = widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
			if isVisible {
				passwordEntry.SetText("")
				showHideButton.SetIcon(theme.VisibilityIcon())
				passwordEntry.Refresh()
			} else {
				passwordEntry.SetText(password["password"])
				showHideButton.SetIcon(theme.VisibilityOffIcon())
				passwordEntry.Refresh()
			}
			isVisible = !isVisible
		})

		folders, err := dbController.GetAllFolders()
		if err != nil {
			log.Println(err)
		}
		var folderKeys []string
		var selectedFolder string
		for folderName, folderID := range folders {
			if fmt.Sprintf("%v", folderID) == password["folder_id"] {
				selectedFolder = folderName
			}
			folderKeys = append(folderKeys, folderName)
		}
		selectFolder := widget.NewSelect(folderKeys, func(selected string) {
			fmt.Println("Selected:", selected)
		})
		selectFolder.Selected = selectedFolder
		selectFolder.Disable()

		encryption, err := dbController.GetAllEncrypted()
		if err != nil {
			log.Fatal(err)
		}
		var encryptionKeys []string
		var selectedEncryption string
		for encryptionName, encryptionID := range encryption {
			if fmt.Sprintf("%v", encryptionID) == password["encrypted_id"] {
				selectedEncryption = encryptionName
			}
			encryptionKeys = append(encryptionKeys, encryptionName)
		}
		selectEncryption := widget.NewSelect(encryptionKeys, func(selected string) { fmt.Println("Selected:", selected) })
		selectEncryption.Selected = selectedEncryption
		selectEncryption.Disable()

		note := widget.NewMultiLineEntry()
		note.Text = password["note"]
		note.Disable()

		var returnButton *widget.Button
		var editButton *widget.Button
		var saveButton *widget.Button
		var cancelButton *widget.Button
		var deleteButton *widget.Button
		returnButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
			controller.ShowScreen(screenReturn, folderName, folderId)
		})

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
							fmt.Printf("Error deleting password: %v\n", err)
						} else {
							fmt.Printf("Deleted password: %v:%v\n", password["label"], password["password"])
						}
						dialog.Hide()
						controller.ShowScreen(screenReturn, folderName, folderId)
					}),
				),
			), w.Canvas())
			dialog.Show()
		})

		editButton = widget.NewButton(localizer.Get("edit"), func() {
			labelEntry.Enable()
			nameEntry.Enable()
			passwordEntry.Enable()
			websiteEntry.Enable()
			selectFolder.Enable()
			selectEncryption.Enable()
			note.Enable()
			editButton.Hide()
			returnButton.Hide()
			saveButton.Show()
			cancelButton.Show()
		})

		saveButton = widget.NewButton(localizer.Get("save"), func() {
			updatedPassword := make(map[string]interface{})

			updatedPassword["label"] = labelEntry.Text
			updatedPassword["name"] = nameEntry.Text
			updatedPassword["website"] = websiteEntry.Text

			if passwordEntry.Text != "" {
				updatedPassword["password"] = passwordEntry.Text
			}
			if _, ok := folders[selectFolder.Selected]; ok {
				updatedPassword["folder_id"] = fmt.Sprintf("%v", folders[selectFolder.Selected])
			}
			if _, ok := folders[selectFolder.Selected]; ok {
				updatedPassword["encrypted_id"] = fmt.Sprintf("%v", encryption[selectEncryption.Selected])
			}
			updatedPassword["note"] = note.Text

			for key, value := range updatedPassword {
				if strValue, ok := value.(string); ok {
					password[key] = strValue
				}
			}

			idStr := password["id"]
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Println("Error al convertir el id a int64:", err)
			}

			err = dbController.EditPassword(id, updatedPassword, contUser.GetCurrentUserPassword())
			if err != nil {
				log.Fatalf(err.Error())
			}

			labelEntry.Disable()
			nameEntry.Disable()
			passwordEntry.Disable()
			websiteEntry.Disable()
			selectFolder.Disable()
			selectEncryption.Disable()
			note.Disable()
			editButton.Show()
			returnButton.Show()
			saveButton.Hide()
			cancelButton.Hide()
		})
		saveButton.Hide()

		cancelButton = widget.NewButton(localizer.Get("cancel"), func() {
			editButton.Show()
			returnButton.Show()
			saveButton.Hide()
			cancelButton.Hide()
		})
		cancelButton.Hide()

		w.SetContent(container.NewVBox(
			container.NewHBox(viewTitle, returnButton, deleteButton),
			passwordLabel,
			labelEntry,
			passwordName,
			nameEntry,
			container.NewHBox(passwordPassword, showHideButton),
			passwordEntry,
			passwordWebsite,
			websiteEntry,
			passwordFolder,
			selectFolder,
			passwordEncryption,
			selectEncryption,
			passwordNote,
			note,
			editButton,
			saveButton,
			cancelButton,
		))
	}
}
