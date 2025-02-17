package ui

import (
	"fmt"
	"log"
	"scribe-nb/scribedb"
	"scribe-nb/note"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	//"github.com/fyne-io/terminal"
)

func OpenNoteWindow(app fyne.App, noteId uint) {
	retreievdNote, err := scribedb.GetNote(noteId)

	if err != nil{
		log.Println("error getting note")
		log.Panic(err)
	}

	noteInfo := note.NoteInfo{
		Id: retreievdNote.Id,
		Notebook: retreievdNote.Notebook,
		DateCreated: retreievdNote.Created,
		DateModified: retreievdNote.Modified,
		Colour: retreievdNote.BackgroundColour,
		Content: retreievdNote.Content,
	}

	if noteInfo.Id != 0{
		noteInfo.NewNote = false
	}else{
		noteInfo.NewNote = true
	}

	if retreievdNote.Pinned > 0{
		noteInfo.Pinned = true
	} else {
		noteInfo.Pinned = false
	}

	//calculate initial note content hash
	note.UpdateHash(&noteInfo)

	noteWindow := app.NewWindow(fmt.Sprintf("Notebook: %s --- Note id: %d", retreievdNote.Notebook, retreievdNote.Id))
	noteWindow.Resize(fyne.NewSize(850, 900))
	entry := widget.NewEntry()
	entry.Text = noteInfo.Content

	cont := container.NewStack(entry)

	noteWindow.SetContent(cont)
	noteWindow.Canvas().Focus(entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", noteInfo.Id))
		fmt.Println("You will want to remove the id from the openNotes list here!!!!!!!")
		noteInfo.Content = entry.Text
		if note.CheckChanges(&retreievdNote,&noteInfo){
			res, err := note.SaveNote(&noteInfo)
			if err != nil{
				log.Println("Error saving note")
				log.Panic()
			}

			if res == 0{
				log.Println("No note was saved (affected rows = 0)")
			}else{
				log.Println("....Note saved successfully....")
			}
		}

	})
	noteWindow.Show()
}
