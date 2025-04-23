package ui

import (
	"scribe-nb/note"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewProperetiesPanel() *fyne.Container {
	themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
	propertiesTitle := widget.NewRichTextFromMarkdown("**Properties**")
	NoteWidgets.propertiesText = widget.NewLabel("")
	vbox := container.NewVBox(propertiesTitle, NoteWidgets.propertiesText)
	propertiesPadded := container.NewPadded(themeBackground, vbox)
	return container.NewStack(propertiesPadded)
}

func ShowProperties(noteInfo *note.NoteInfo) {
	if NoteContainers.propertiesPanel.Hidden {
		text := note.GetPropertiesText(noteInfo)
		NoteWidgets.propertiesText.SetText(text)
		NoteContainers.propertiesPanel.Show()
	} else {
		NoteContainers.propertiesPanel.Hide()
		AppContainers.singleNoteStack.Refresh() //only needed for single page view - IMPROVE THIS!!!!!!!
	}
}

func UpdateProperties(noteInfo *note.NoteInfo) {
	if !NoteContainers.propertiesPanel.Hidden {
		text := note.GetPropertiesText(noteInfo)
		NoteWidgets.propertiesText.SetText(text)
		NoteWidgets.propertiesText.Refresh()
	}
}
