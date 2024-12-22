package controllers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Screen es una función que configura el contenido de una pantalla en la ventana principal,
// y acepta parámetros opcionales.
type Screen func(w fyne.Window, params ...interface{})

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
// También puede pasar parámetros opcionales para personalizar el contenido.
func (c *ControllerScreen) ShowScreen(name string, params ...interface{}) {
	if screen, exists := c.screens[name]; exists {
		screen(c.window, params...)
	} else {
		c.window.SetContent(container.NewCenter(
			widget.NewLabel("Error: Screen not found"),
		))
	}
}
