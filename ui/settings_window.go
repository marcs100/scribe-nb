package ui

import (
	"fmt"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsWindow() {
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
	//notesLimitHBox := container.NewHBox(recentNotesLimitLabel, recentNotesLimitEntry)
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
	gridLimitStack := container.NewStack(gridLimitEntry)
	gridLimitGrid := container.NewGridWithRows(1, gridLimitLabel, gridLimitStack)
	gridLimitEntry.SetText(fmt.Sprintf("%d", Conf.Settings.GridMaxPages))

	vbox := container.NewVBox(viewHeading, viewGrid, notesLimitGrid, layoutHeading, layoutGrid, gridLimitGrid)

	stack := container.NewStack(bg, vbox)

	settingsWindow.SetContent(stack)
	settingsWindow.Show()
}
