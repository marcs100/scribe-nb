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
	modDarkColour, _ := RGBStringToFyneColor(Conf.Settings.DarkColour)
	modLightColour, _ := RGBStringToFyneColor(Conf.Settings.LightColour)
	switch themeVariant {
	case theme.VariantDark:
		themeColour = modDarkColour
	case theme.VariantLight:
		themeColour = modLightColour
	}
	return themeColour
}
