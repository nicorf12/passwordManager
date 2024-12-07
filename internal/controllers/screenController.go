package controllers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Screen es una función que configura el contenido de una pantalla en la ventana principal.
type Screen func(w fyne.Window)

// ControllerScreen gestiona los cambios de pantalla.
type ControllerScreen struct {
	window  fyne.Window
	screens map[string]Screen
}

// NewControllerScreen crea e inicializa un ControllerScreen.
func NewControllerScreen(w fyne.Window) *ControllerScreen {
	return &ControllerScreen{
		window:  w,
		screens: make(map[string]Screen),
	}
}

// RegisterScreen registra una pantalla con un nombre único.
func (c *ControllerScreen) RegisterScreen(name string, screen Screen) {
	c.screens[name] = screen
}

// ShowScreen cambia a la pantalla registrada con el nombre proporcionado.
func (c *ControllerScreen) ShowScreen(name string) {
	if screen, exists := c.screens[name]; exists {
		screen(c.window)
	} else {
		c.window.SetContent(container.NewCenter(
			widget.NewLabel("Error: Screen not found"),
		))
	}
}
