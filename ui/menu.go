package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"password_manager/internal/controllers"
	"password_manager/localization"
	"sort"
)

func menu(controller *controllers.ControllerScreen, contUser *controllers.ControllerUser, dbController *controllers.DBController, localizer *localization.Localizer) *fyne.Container {
	menuTitle := widget.NewLabel(localizer.Get("menu"))
	logoutButton := widget.NewButton(localizer.Get("logout"), func() {
		contUser.Logout()
		controller.ShowScreen("login")
	})
	optionsTitle := widget.NewLabel(localizer.Get("options"))
	allButton := widget.NewButtonWithIcon(localizer.Get("all"), theme.MediaPlayIcon(), func() {
		controller.ShowScreen("main")
	})
	favoritesButton := widget.NewButtonWithIcon(localizer.Get("favorites"), theme.CheckButtonCheckedIcon(), func() {
		controller.ShowScreen("favorites")
	})
	categoriesTitle := widget.NewLabel(localizer.Get("categories"))

	categories := container.NewVBox()
	folders, err := dbController.GetAllFolders()
	if err != nil {
		log.Println("Error fetching folders:", err)
	} else {
		var folderKeys []string
		for name := range folders {
			folderKeys = append(folderKeys, name)
		}

		sort.Strings(folderKeys)

		for _, name := range folderKeys {
			folderButton := widget.NewButtonWithIcon(name, theme.FolderIcon(), func() {
				controller.ShowScreen("folder", name, folders[name])
			})
			categories.Add(folderButton)
		}
	}

	var (
		addCategoryButton *widget.Button
		confirmAddButton  *widget.Button
		cancelAddButton   *widget.Button
	)

	folderNameEntry := widget.NewEntry()
	folderNameEntry.SetPlaceHolder(localizer.Get("folderName"))
	folderNameEntry.Hide()

	confirmAddButton = widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
		_, err := dbController.InsertFolder(folderNameEntry.Text)
		if err != nil {
			return
		}

		folders, err := dbController.GetAllFolders()
		if err != nil {
			log.Println("Error fetching folders:", err)
			return
		}

		categories.Objects = nil
		var folderKeys []string
		for name := range folders {
			folderKeys = append(folderKeys, name)
		}

		sort.Strings(folderKeys)

		for _, name := range folderKeys {
			folderButton := widget.NewButtonWithIcon(name, theme.FolderIcon(), func() {
				controller.ShowScreen("folder", name, folders[name])
			})
			categories.Add(folderButton)
		}

		categories.Refresh()

		addCategoryButton.Show()
		confirmAddButton.Hide()
		cancelAddButton.Hide()
		folderNameEntry.Hide()
		folderNameEntry.SetText("")
	})

	confirmAddButton.Hide()

	cancelAddButton = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		addCategoryButton.Show()
		confirmAddButton.Hide()
		cancelAddButton.Hide()
		folderNameEntry.Hide()
		folderNameEntry.Text = ""
	})
	cancelAddButton.Hide()

	addCategoryButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		addCategoryButton.Hide()
		confirmAddButton.Show()
		cancelAddButton.Show()
		folderNameEntry.Show()
	})

	menu := container.NewVBox(
		menuTitle,
		logoutButton,
		canvas.NewLine(color.White),
		optionsTitle,
		allButton,
		favoritesButton,
		canvas.NewLine(color.White),
		categoriesTitle,
		categories,
		folderNameEntry,
		confirmAddButton,
		cancelAddButton,
		addCategoryButton,
	)

	return menu
}
