package ui

import (
	"fyne.io/fyne/v2/theme"
	"image/color"
)

// Background colour for notes based on current theme variane light/dark
func GetThemeColour() color.Color {
	var themeColour color.Color
	themeVariant := mainApp.Settings().ThemeVariant()
	themeColour = mainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
	modDarkColour, _ := RGBStringToFyneColor("#2f2f2f")
	modLightColour, _ := RGBStringToFyneColor("#e2e2e2")
	switch themeVariant {
	case theme.VariantDark:
		themeColour = modDarkColour
	case theme.VariantLight:
		themeColour = modLightColour
	}
	return themeColour
}
