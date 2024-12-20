package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"password_manager/localization"
)

func windowImportExport(app fyne.App, f func(string, string, string), localizer *localization.Localizer, text string, errorLabel *canvas.Text) fyne.Window {
	w := app.NewWindow(localizer.Get(text))
	w.Resize(fyne.NewSize(600, 600))
	icono, err := fyne.LoadResourceFromPath("resources/dragon.png")
	if err != nil {
		log.Printf("Error cargando el ícono: %v", err)
		icono = theme.ErrorIcon()
	}
	w.SetIcon(icono)

	q0 := widget.NewLabel(localizer.Get("question0"))
	q1 := widget.NewLabel(localizer.Get("question1"))
	q2 := widget.NewLabel(localizer.Get("question2"))

	answer0 := widget.NewEntry()
	answer1 := widget.NewEntry()
	answer2 := widget.NewEntry()

	button := widget.NewButton(localizer.Get(text), func() {
		if answer0.Text == "" || answer1.Text == "" || answer2.Text == "" {
			fmt.Println("any incomplete answer")
			return
		}
		f(answer0.Text, answer1.Text, answer2.Text)
	})

	content := container.NewVBox(
		container.NewVBox(q0, answer0),
		container.NewVBox(q1, answer1),
		container.NewVBox(q2, answer2),
		button,
		errorLabel,
	)
	w.SetContent(content)
	return w
}
