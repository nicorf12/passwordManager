package themes

import (
	"fyne.io/fyne/v2"
	"image/color"
)

var (
	pink                color.Color = color.RGBA{R: 255, G: 20, B: 147, A: 255}
	lightPink           color.Color = color.RGBA{R: 255, G: 240, B: 245, A: 255}
	darkPink            color.Color = color.RGBA{R: 229, G: 44, B: 100, A: 255}
	white               color.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	softButtonColorPink color.Color = color.RGBA{R: 248, G: 205, B: 248, A: 255}
	blue                color.Color = color.RGBA{R: 70, G: 130, B: 180, A: 255}
	lightBlue           color.Color = color.RGBA{R: 173, G: 216, B: 230, A: 255}
	darkBlue            color.Color = color.RGBA{R: 6, G: 2, B: 112, A: 255}
	softButtonColorBlue color.Color = color.RGBA{R: 135, G: 206, B: 250, A: 255}
	darkBackground      color.Color = color.RGBA{R: 18, G: 18, B: 18, A: 255}
	lightForeground     color.Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	buttonBlue          color.Color = color.RGBA{R: 50, G: 150, B: 255, A: 255}
	disabledLightBlue   color.Color = color.RGBA{R: 100, G: 170, B: 255, A: 255}
	transparent                     = color.Transparent
)

func GetTheme(theme string) fyne.Theme {
	switch theme {
	case "Dark":
		return CustomDarkTheme{}
	case "Light":
		return CustomLightTheme{}
	case "Pink":
		return CustomPinkTheme{}
	case "Blue":
		return CustomBlueTheme{}
	default:
		return CustomDarkTheme{}
	}
}
