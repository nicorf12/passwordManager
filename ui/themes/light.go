package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// CustomLightTheme implementa fyne.Theme
type CustomLightTheme struct{}

// Color devuelve los colores personalizados para el tema
func (c CustomLightTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	case theme.ColorNameForeground:
		return color.Black
	case theme.ColorNameButton:
		return transparent
	case theme.ColorNamePrimary:

		return softButtonColorBlue
	default:
		return theme.LightTheme().Color(name, variant)
	}
}

// Font devuelve la fuente personalizada para el tema
func (c CustomLightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon devuelve el icono personalizado para el tema
func (c CustomLightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size devuelve el tamaño de los elementos en el tema
func (c CustomLightTheme) Size(name fyne.ThemeSizeName) float32 {
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
