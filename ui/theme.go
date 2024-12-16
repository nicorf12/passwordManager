package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type customDarkTheme struct{}

func (c customDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameDisabled:
		return color.RGBA{R: 100, G: 170, B: 255, A: 255}
	case theme.ColorNameButton:
		return color.Transparent //color.RGBA{R: 50, G: 150, B: 255, A: 255}
	default:
		return theme.DarkTheme().Color(name, variant)
	}
}

func (c customDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

func (c customDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (c customDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DarkTheme().Size(name)
}

type customTheme struct {
	fyne.Theme
}

func (t *customTheme) ButtonColor() color.Color {
	return color.Transparent // Fondo transparente
}

func (t *customTheme) DisabledButtonColor() color.Color {
	return color.Transparent // Fondo transparente para botones deshabilitados
}
