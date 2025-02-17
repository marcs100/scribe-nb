package ui

import (
	"fmt"
	"log"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	//"github.com/fyne-io/terminal"
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
	entry.Text = note.Content

	/*term := terminal.New()
	go func(){
		_ = term.RunLocalShell()
	}() */

	cont := container.NewStack(entry)

	noteWindow.SetContent(cont)
	noteWindow.Show()
}


