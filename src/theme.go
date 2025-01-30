package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type forcedVariant struct {
	fyne.Theme

	variant fyne.ThemeVariant
}

var (
	PrimaryGreenColor    = color.RGBA{R: 141, G: 181, B: 68, A: 255}
	SecondaryYellowColor = color.RGBA{R: 206, G: 187, B: 91, A: 255}
)

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	/*if name == theme.ColorNameBackground {
		return color.White
	}*/
	switch name {
	case theme.ColorNamePrimary:
		return PrimaryGreenColor
	case theme.ColorNameWarning:
		return SecondaryYellowColor
	case theme.ColorNameFocus:
		return color.RGBA{122, 181, 14, 35}
	case theme.ColorNameSelection:
		return color.RGBA{122, 150, 74, 255}

	default:
		return f.Theme.Color(name, f.variant)
	}
}

func (f *forcedVariant) Size(s fyne.ThemeSizeName) float32 {
	if s == theme.SizeNameSelectionRadius {
		return 10
	} else if s == theme.SizeNameSeparatorThickness {
		return 2
	}
	return f.Theme.Size(s)
}

/*func (f *forcedVariant) Font(style fyne.TextStyle) fyne.Resource {
	return f.Theme.Font(fyne.TextStyle{Bold: true})
}*/
