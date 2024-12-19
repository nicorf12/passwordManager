package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// CustomBlueTheme implementa fyne.Theme
type CustomBlueTheme struct{}

// Color devuelve los colores personalizados para el tema
func (c CustomBlueTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return lightBlue
	case theme.ColorNameForeground:
		return blue
	case theme.ColorNameButton:
		return softButtonColorBlue
	case theme.ColorNameInputBackground:
		return white
	case theme.ColorNameDisabled:
		return darkBlue
	case theme.ColorNameError:
		return darkBlue
	case theme.ColorNameFocus:
		return white
	case theme.ColorNameHover:
		return color.RGBA{R: 119, G: 115, B: 236, A: 255}
	case theme.ColorNameInputBorder:
		return blue
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 119, G: 115, B: 236, A: 255}
	default:
		return theme.LightTheme().Color(name, variant)
	}
}

// Font devuelve la fuente personalizada para el tema
func (c CustomBlueTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon devuelve el icono personalizado para el tema
func (c CustomBlueTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size devuelve el tamaño de los elementos en el tema
func (c CustomBlueTheme) Size(name fyne.ThemeSizeName) float32 {
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
