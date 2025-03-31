package ui

import (
	"log"

	"fyne.io/fyne/v2/theme"
)

// Background colour for notes based on current theme variane light/dark
func GetThemeColours() AppColours {
	var appColours AppColours
	var err error
	themeVariant := mainApp.Settings().ThemeVariant()
	appColours.MainBgColour = mainApp.Settings().Theme().Color(theme.ColorNameBackground, themeVariant)
	appColours.NoteBgColour = appColours.MainBgColour
	switch themeVariant {
	case theme.VariantDark:
		appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourBg)
		if err != nil {
			log.Println("************* here we are **********************************")
			log.Println(Conf.Settings.DarkColourBg)
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.DarkColourNote)
		if err != nil {
			log.Panicln(err)
		}
	case theme.VariantLight:
		appColours.MainBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourBg)
		if err != nil {
			log.Panicln(err)
		}
		appColours.NoteBgColour, err = RGBStringToFyneColor(Conf.Settings.LightColourNote)
		if err != nil {
			log.Panicln(err)
		}
	}
	return appColours
}
