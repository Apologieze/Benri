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

func (f *forcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	/*if name == theme.ColorNameBackground {
		return color.White
	}*/
	theme.SelectionRadiusSize()
	return f.Theme.Color(name, f.variant)
}

func (f *forcedVariant) Size(s fyne.ThemeSizeName) float32 {
	if s == theme.SizeNameSelectionRadius {
		return 10
	} else if s == theme.SizeNameSeparatorThickness {
		return 2
	}
	return f.Theme.Size(s)
}
