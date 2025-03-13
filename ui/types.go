package ui

import (
	"fmt"
	"image/color"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PageViewStatus struct{
	NumberOfPages int
	CurrentPage int
}

func (pv *PageViewStatus) PageForward() int{
	if pv.CurrentPage < pv.NumberOfPages{
		pv.CurrentPage += 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageViewStatus) PageBack() int{
	if pv.CurrentPage > 1{
		pv.CurrentPage -= 1
		return pv.CurrentPage
	}
	return -1
}

func (pv *PageViewStatus) Reset(){
	pv.CurrentPage = 0
	pv.NumberOfPages = 0
}

func (pv *PageViewStatus) GetLabel()string{
	return fmt.Sprintf("Page: %d of %d",pv.CurrentPage,pv.NumberOfPages)
}


type ApplicationContainers struct{
	grid *fyne.Container
	singleNoteStack *fyne.Container
	mainGridContainer *container.Scroll
	mainPageContainer *container.Scroll
}

type ApplicationWidgets struct{
	toolbar *widget.Toolbar
	singleNotePage *widget.RichText
	viewLabel *widget.Label
	pageLabel *widget.Label
	notebooksList *widget.List
}

type ApplicationStatus struct{
	currentView string
	currentNotebook string
	currentLayout string
	notes []scribedb.NoteData
	notebooks []string
	themeBgColour color.Color
	openNotes []uint //maintain a list of notes that are currently open
	noteSize fyne.Size
}

