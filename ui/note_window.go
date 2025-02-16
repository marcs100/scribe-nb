package ui

import (
	"scribe-nb/scribedb"
	"log"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func OpenNoteWindow(app fyne.App, noteId uint){

	note,err := scribedb.GetNote(noteId)

	if err != nil{
		log.Printf("Error getting note %d", noteId)
		return
	}

	noteWindow := app.NewWindow(fmt.Sprintf("Notebook: %s --- Note id: %d",note.Notebook,note.Id))
	noteWindow.Resize(fyne.NewSize(850,900))
	entry := widget.NewEntry()
	cont := container.NewStack(entry)

	noteWindow.SetContent(cont)
	noteWindow.Show()

}


