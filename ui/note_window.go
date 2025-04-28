package ui

import (
	"fmt"
	"log"
	"scribe-nb/note"
	"scribe-nb/scribedb"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	//"fyne.io/fyne/v2/data/binding"
	//"github.com/fyne-io/terminal"
)

func OpenNoteWindow(noteId uint) {
	var noteInfo note.NoteInfo
	var retrievedNote scribedb.NoteData
	var err error

	if noteId != 0 {
		retrievedNote, err = scribedb.GetNote(noteId)
		if err != nil {
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	}

	noteWindow := mainApp.NewWindow("")
	noteContainer := NewNoteContainer(noteId, &noteInfo, &retrievedNote, noteWindow)
	//fmt.Println(fmt.Sprintf("************Notebook is %s", "debug"))
	noteWindow.SetTitle(fmt.Sprintf("Notebook: %s", noteInfo.Notebook))
	noteWindow.Resize(fyne.NewSize(900, 750))

	noteWindow.SetContent(noteContainer)
	noteWindow.Canvas().Focus(NoteWidgets.entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", noteInfo.Id))
		SaveNote(&noteInfo, &retrievedNote)

		if index := slices.Index(AppStatus.openNotes, noteInfo.Id); index != -1 {
			AppStatus.openNotes = slices.Delete(AppStatus.openNotes, index, index+1)
		}
	})

	AddNoteKeyboardShortcuts(&noteInfo, noteWindow)

	if noteInfo.NewNote {
		SetEditMode(noteWindow)
	}

	noteWindow.Show()
}
