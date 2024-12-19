package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// CustomDarkTheme implementa fyne.Theme
type CustomDarkTheme struct{}

// Color devuelve los colores personalizados para el tema
func (c CustomDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return darkBackground
	case theme.ColorNameForeground:
		return lightForeground
	case theme.ColorNameButton:
		return buttonBlue
	case theme.ColorNameDisabled:
		return disabledLightBlue
	case theme.ColorNamePrimary:

		return softButtonColorBlue
	default:
		return theme.DarkTheme().Color(name, variant)
	}
}

// Font devuelve la fuente personalizada para el tema
func (c CustomDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DarkTheme().Font(style)
}

// Icon devuelve el icono personalizado para el tema
func (c CustomDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

// Size devuelve el tamaño de los elementos en el tema
func (c CustomDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 16
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBar:
		return 6
	default:
		return theme.LightTheme().Size(name)
	}
}
