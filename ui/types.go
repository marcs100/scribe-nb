package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"scribe-nb/scribedb"
	"image/color"
)

type PageView struct{
	NumberOfPages int
	CurrentPage int
}

func (pv *PageView) PageForward() int{
	if pv.CurrentPage < pv.NumberOfPages{
		pv.CurrentPage += 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageView) PageBack() int{
	if pv.CurrentPage > 1{
		pv.CurrentPage -= 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageView) Reset(){
	pv.CurrentPage = 0
	pv.NumberOfPages = 0
}


type AppContainers struct{
	grid *fyne.Container
	singleNoteStack *fyne.Container
	mainGridContainer *container.Scroll
	mainPageContainer *container.Scroll

}

type AppWidgets struct{
	toolbar *widget.Toolbar
	singleNotePage *widget.RichText
}

type AppStatus struct{
	currentView string
	currentNotebook string
	currentLayout string
	notes []scribedb.NoteData
	themeBgColour color.Color
	openNotes []uint //maintain a list of notes that are currently open
}

