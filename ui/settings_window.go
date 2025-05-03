package ui

import (
	"fmt"
	"scribe-nb/config"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsWindow() {
	
	origConf := CopySettings()
	
	var themeVar theme_variant
	switch Conf.Settings.ThemeVariant {
	case "light":
		themeVar = LIGHT_THEME
	case "dark":
		themeVar = DARK_THEME
	case "system":
		themeVar = SYSTEM_THEME
	}
	AppTheme = GetThemeColours(themeVar)

	settingsWindow := mainApp.NewWindow("Settings")
	settingsWindow.Resize(fyne.NewSize(500, 400))

	bg := canvas.NewRectangle(AppTheme.MainBgColour)

	viewHeading := widget.NewRichTextFromMarkdown("### View")
	viewLabel := widget.NewLabel("  Default View:          ")
	viewSelect := widget.NewSelect([]string{"pinned", "recent"}, func(sel string) {
		Conf.Settings.InitialView = sel
	})
	viewSelect.SetSelected(Conf.Settings.InitialView)
	viewGrid := container.NewGridWithRows(1, viewLabel, viewSelect)

	recentNotesLimitLabel := widget.NewLabel("  Recent Note Limit:")
	recentNotesLimitEntry := widget.NewEntry()
	recentNotesLimitEntry.SetText(fmt.Sprintf("%d", Conf.Settings.RecentNotesLimit))
	recentNotesLimitEntry.OnChanged = func(input string) {
		_, err := strconv.Atoi(input)
		if err != nil {
			recentNotesLimitEntry.SetText("")
		}
	}

	notesLimitGrid := container.NewGridWithRows(1, recentNotesLimitLabel, recentNotesLimitEntry)

	layoutHeading := widget.NewRichTextFromMarkdown("### Layout")
	layoutLabel := widget.NewLabel("  Default Layout:")
	layoutSelect := widget.NewSelect([]string{"grid", "page"}, func(sel string) {
		Conf.Settings.InitialLayout = sel
	})
	layoutSelect.Selected = Conf.Settings.InitialLayout
	layoutGrid := container.NewGridWithRows(1, layoutLabel, layoutSelect)

	gridLimitLabel := widget.NewLabel("  Notes per Page Limit:")
	gridLimitEntry := widget.NewEntry()
	gridLimitEntry.OnChanged = func(input string) {
		_, err := strconv.Atoi(input)
		if err != nil {
			gridLimitEntry.SetText("")
		}
	}
	gridLimitStack := container.NewStack(gridLimitEntry)
	gridLimitGrid := container.NewGridWithRows(1, gridLimitLabel, gridLimitStack)
	gridLimitEntry.SetText(fmt.Sprintf("%d", Conf.Settings.GridMaxPages))

	appearanceHeading := widget.NewRichTextFromMarkdown("### Appearance")
	appearanceLabel := widget.NewLabel("  Theme:")
	appearanceSelect := widget.NewSelect([]string{"light", "dark", "system"}, func(sel string) {
		Conf.Settings.ThemeVariant = sel
	})
	appearanceSelect.Selected = Conf.Settings.ThemeVariant
	appearanceGrid := container.NewGridWithRows(1, appearanceLabel, appearanceSelect)

	vbox := container.NewVBox(
		viewHeading,
		viewGrid,
		notesLimitGrid,
		layoutHeading,
		layoutGrid,
		gridLimitGrid,
		appearanceHeading,
		appearanceGrid)

	stack := container.NewStack(bg, vbox)

	settingsWindow.SetContent(stack)
	settingsWindow.Show()
}

func CopySettings() config.Config{
	return config.Config{
		Title: Conf.Title,
		Settings: config.AppSettings{
			Database: Conf.Settings.Database,
			RecentNotesLimit: Conf.Settings.RecentNotesLimit,
			NoteWidth: Conf.Settings.NoteWidth,
			NoteHeight: Conf.Settings.NoteHeight,
			InitialView: Conf.Settings.InitialView,
			InitialLayout: Conf.Settings.InitialLayout,
			GridMaxPages: Conf.Settings.GridMaxPages,
			ThemeVariant: Conf.Settings.ThemeVariant,
			DarkColourNote: Conf.Settings.DarkColourNote,
			LightColourNote: Conf.Settings.LightColourNote,
			DarkColourBg: Conf.Settings.DarkColourBg,
			LightColourBg: Conf.Settings.LightColourBg,
		},
	}
}
