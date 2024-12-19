package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// CustomPinkTheme implementa fyne.Theme
type CustomPinkTheme struct{}

// Color devuelve los colores personalizados para el tema
func (c CustomPinkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return lightPink
	case theme.ColorNameForeground:
		return pink
	case theme.ColorNameButton:
		return softButtonColorPink
	case theme.ColorNameInputBackground:
		return white
	case theme.ColorNameDisabled:
		return darkPink
	case theme.ColorNameError:
		return darkPink
	case theme.ColorNameFocus:
		return white
	case theme.ColorNameHover:
		return color.RGBA{R: 255, G: 182, B: 193, A: 255}
	case theme.ColorNameInputBorder:
		return pink
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 219, G: 112, B: 147, A: 255}
	default:
		return theme.LightTheme().Color(name, variant)
	}
}

// Font devuelve la fuente personalizada para el tema
func (c CustomPinkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon devuelve el icono personalizado para el tema
func (c CustomPinkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size devuelve el tamaño de los elementos en el tema
func (c CustomPinkTheme) Size(name fyne.ThemeSizeName) float32 {
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
