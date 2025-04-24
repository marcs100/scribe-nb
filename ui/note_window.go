package ui

import (
	"fmt"
	"log"
	"scribe-nb/note"
	"scribe-nb/scribedb"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
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

func AddNoteKeyboardShortcuts(noteInfo *note.NoteInfo, parentWindow fyne.Window) {
	//Keyboard shortcut to set edit mode
	ctrl_e := &desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierControl,
	}

	parentWindow.Canvas().AddShortcut(ctrl_e, func(shortcut fyne.Shortcut) {
		SetEditMode(parentWindow)
	})

	//Keyboard shortcut to set view mode
	ctrl_q := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}

	parentWindow.Canvas().AddShortcut(ctrl_q, func(shortcut fyne.Shortcut) {
		SetViewMode(parentWindow)
	})

	//Keyboard shortcut to pin/unpin notes
	ctrl_p := &desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl,
	}

	parentWindow.Canvas().AddShortcut(ctrl_p, func(shortcut fyne.Shortcut) {
		PinNote(noteInfo)
	})

	//Keyboard shortcut to change note colour
	ctrl_h := &desktop.CustomShortcut{
		KeyName:  fyne.KeyH,
		Modifier: fyne.KeyModifierControl,
	}

	parentWindow.Canvas().AddShortcut(ctrl_h, func(shortcut fyne.Shortcut) {
		ChangeNoteColour(noteInfo, parentWindow)
	})
}
