package ui

import (
	"fmt"
	"image/color"
	"log"
	"scribe-nb/note"
	"scribe-nb/scribedb"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/data/binding"
	//"github.com/fyne-io/terminal"
)

func OpenNoteWindow(noteId uint) {
	var err error
	var retrievedNote scribedb.NoteData
	var noteInfo note.NoteInfo

	if noteId == 0 {
		//New note
		noteInfo = note.NoteInfo{
			Id:           noteId,
			Notebook:     "General",
			DateCreated:  "",
			DateModified: "",
			Pinned:       false,
			PinnedDate:   "",
			Colour:       "#FFFFFF",
			Content:      "",
			Deleted:      false,
		}
	} else {
		//existing note
		retrievedNote, err = scribedb.GetNote(noteId)

		if err != nil {
			log.Println("error getting note")
			dialog.ShowError(err, noteWindow)
			log.Panic(err)
		}

		noteInfo = note.NoteInfo{
			Id:           retrievedNote.Id,
			Notebook:     retrievedNote.Notebook,
			DateCreated:  retrievedNote.Created,
			DateModified: retrievedNote.Modified,
			PinnedDate:   retrievedNote.PinnedDate,
			Colour:       retrievedNote.BackgroundColour,
			Content:      retrievedNote.Content,
			Deleted:      false,
		}

		if retrievedNote.Pinned > 0 {
			noteInfo.Pinned = true
		} else {
			noteInfo.Pinned = false
		}
	}

	if noteInfo.Id != 0 {
		noteInfo.NewNote = false
	} else {
		noteInfo.NewNote = true
	}

	//calculate initial note content hash
	note.UpdateHash(&noteInfo)

	noteWindow = mainApp.NewWindow(fmt.Sprintf("Notebook: %s", noteInfo.Notebook))
	noteWindow.Resize(fyne.NewSize(900, 750))

	//NoteWidgets.entry = widget.NewMultiLineEntry()
	ctrl_q := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}

	NoteWidgets.entry = NewEntryCustom(ctrl_q, func() {
		SetViewMode()
	})

	//NoteWidgets.entry = widget.NewMultiLineEntry()
	NoteWidgets.entry.Text = noteInfo.Content
	NoteWidgets.entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
	noteColour, _ := RGBStringToFyneColor(noteInfo.Colour)

	NoteCanvas.noteBackground = canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF" || noteInfo.Colour == "#000000" {
		NoteCanvas.noteBackground = canvas.NewRectangle(AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(NoteCanvas.noteBackground)

	NoteWidgets.markdownText = widget.NewRichTextFromMarkdown(noteInfo.Content)
	NoteWidgets.markdownText.Wrapping = fyne.TextWrapWord
	NoteWidgets.markdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, NoteWidgets.markdownText)
	NoteContainers.markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(NoteContainers.markdown)
	background := canvas.NewRectangle(AppTheme.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, NoteWidgets.entry)

	var win *fyne.Container

	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if noteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	NoteWidgets.pinButton = widget.NewButtonWithIcon("", btnIcon, func() {
		PinNote(&noteInfo)
	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	changeNotebookBtn := NewChangeNotebookButton(&noteInfo)

	colourButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		ChangeNoteColour(&noteInfo)
	})

	NoteWidgets.deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteNote(&noteInfo)
	})

	propertiesButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() { ShowProperties(&noteInfo) })

	NoteWidgets.deleteButton.Hide()

	NoteWidgets.modeSelect = widget.NewRadioGroup([]string{EDIT_MODE, VIEW_MODE}, func(value string) {
		switch value {
		case EDIT_MODE:
			SetEditMode()
		case VIEW_MODE:
			SetViewMode()
		}
	})

	//Experimenting with preoperties side panel for notes
	propertiesTitle := widget.NewRichTextFromMarkdown("**Properties**")
	NoteWidgets.propertiesText = widget.NewLabel("")
	vbox := container.NewVBox(propertiesTitle, NoteWidgets.propertiesText)
	propertiesPadded := container.NewPadded(themeBackground, vbox)
	NoteContainers.propertiesPanel = container.NewStack(propertiesPadded)
	//*******************************************************

	NoteWidgets.modeSelect.SetSelected("View")
	NoteWidgets.modeSelect.Horizontal = true
	toolbar := container.NewHBox(NoteWidgets.modeSelect, spacerLabel, NoteWidgets.pinButton, colourButton, changeNotebookBtn, propertiesButton, NoteWidgets.deleteButton)
	win = container.NewBorder(toolbar, nil, nil, NoteContainers.propertiesPanel, content)

	NoteContainers.propertiesPanel.Hide()

	noteWindow.SetContent(win)
	noteWindow.Canvas().Focus(NoteWidgets.entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", noteInfo.Id))
		noteInfo.Content = NoteWidgets.entry.Text
		var noteChanges note.NoteChanges
		if noteInfo.NewNote {
			if noteInfo.Content != "" {
				noteChanges.ContentChanged = true
			}
		} else {
			noteChanges = note.CheckChanges(&retrievedNote, &noteInfo)
		}
		//if contentChanged{
		if noteChanges.ContentChanged || noteChanges.ParamsChanged {
			res, err := note.SaveNote(&noteInfo)
			if err != nil {
				log.Println("Error saving note")
				dialog.ShowError(err, mainWindow)
				//log.Panic()
			}

			if res == 0 {
				log.Println("No note was saved (affected rows = 0)")
			} else {
				log.Println("....Note updates successfully....")
				go UpdateView()
			}
		} else if noteChanges.PinStatusChanged {
			// we do not want a create or modified time stamp for just pinning/unpinning notes
			res, err := note.SaveNoteNoTimeStamp(&noteInfo)
			if err != nil {
				log.Println("Error saving note")
				dialog.ShowError(err, mainWindow)
				//log.Panic()
			}

			if res == 0 {
				log.Println("No note was saved (affected rows = 0)")
			} else {
				log.Println("....Note updates successfully....")
				go UpdateView()
			}
		}

		if index := slices.Index(AppStatus.openNotes, noteInfo.Id); index != -1 {
			AppStatus.openNotes = slices.Delete(AppStatus.openNotes, index, index+1)
		}

	})

	AddNoteKeyboardShortcuts(&noteInfo)

	if noteInfo.NewNote {
		SetEditMode()
	}

	noteWindow.Show()
}

func NewChangeNotebookButton(noteInfo *note.NoteInfo) *widget.Button {
	changeNotebookBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		var notebooks []string
		var err error
		if notebooks, err = scribedb.GetNotebooks(); err != nil {
			log.Println("Error getting notebook")
			dialog.ShowError(err, mainWindow)
			log.Panicln(err)
		}
		nbMenu := fyne.NewMenu("Select Notebook")

		//Add new notebook entry to menu
		nbMenuItem := fyne.NewMenuItem("*New*", func() {
			//fmt.Println("Need to ask use for new notebook name here!!!!!!!!")
			notebookEntry := widget.NewEntry()
			eNotebookEntry := widget.NewFormItem("Name", notebookEntry)
			newNotebookDialog := dialog.NewForm("New Notebook?", "OK", "Cancel", []*widget.FormItem{eNotebookEntry}, func(confirmed bool) {
				if confirmed {
					//check that the notebook does not already exist
					exists, err := scribedb.CheckNotebookExists(notebookEntry.Text)
					if err == nil {
						if exists == false {
							//chnage notebook to this new notebook
							noteInfo.Notebook = notebookEntry.Text
							noteWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
							_, err = note.SaveNote(noteInfo)
							if err != nil {
								log.Print("Error saving note: ")
								dialog.ShowError(err, mainWindow)
								//log.Panic(err)
							}
							UpdateNotebooksList()
							UpdateProperties(noteInfo)
						}
					} else {
						dialog.ShowError(err, mainWindow)
						log.Panicln(fmt.Sprintf("Error check notebook exists: %s", err))
					}
				}
			}, noteWindow)
			newNotebookDialog.Show()
		})

		nbMenu.Items = append(nbMenu.Items, nbMenuItem)

		//Now add all the existing notebooks to the menu
		for _, notebook := range notebooks {
			menuItem := fyne.NewMenuItem(notebook, func() {
				noteInfo.Notebook = notebook
				//fmt.Println("Change notebook to " + notebook)
				noteWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
				UpdateProperties(noteInfo)
			})
			nbMenu.Items = append(nbMenu.Items, menuItem)
		}

		popUpMenu := widget.NewPopUpMenu(nbMenu, noteWindow.Canvas())
		//popUpMenu.Show()
		pos := fyne.NewPos(250, 40)
		popUpMenu.ShowAtPosition(pos)
		//popUpMenu.ShowAtPosition(e.Position.AddXY(150,0))

	})

	return changeNotebookBtn
}

func DeleteNote(noteInfo *note.NoteInfo) {
	dialog.ShowConfirm("Delete note", "Are you sure?", func(confirm bool) {
		if confirm {
			var res int64
			var err error = nil
			if noteInfo.NewNote {
				res = 1
			} else {
				res, err = scribedb.DeleteNote(noteInfo.Id)
			}

			if res == 0 || err != nil {
				log.Println("Error deleting note - panic!")
				dialog.ShowError(err, mainWindow)
				//log.Panicln(err)
			} else {
				noteInfo.Deleted = true
				noteWindow.Close()
			}
		}
	}, noteWindow)
}

func PinNote(noteInfo *note.NoteInfo) {
	var res int64
	var err error = nil
	if noteInfo.Pinned {
		if noteInfo.NewNote {
			//new note that hasn't been saved yet'
			noteInfo.Pinned = false
			res = 1
		} else {
			res, err = scribedb.UnpinNote(noteInfo.Id)
		}

		if err == nil && res == 1 {
			noteInfo.Pinned = false
			noteInfo.PinnedDate = ""
			if NoteWidgets.pinButton != nil {
				NoteWidgets.pinButton.SetIcon(theme.RadioButtonIcon())
				NoteWidgets.pinButton.Refresh()
			}
		}
	} else {
		if noteInfo.Id == 0 {
			//new note that hasn't been saved yet'
			noteInfo.Pinned = true
			res = 1
		} else {
			res, err = scribedb.PinNote(noteInfo.Id)
			pinnedDate, err := scribedb.GetPinnedDate(int(noteInfo.Id))
			if err != nil {
				log.Println(fmt.Sprintf("Error getting pinned date: &s", err))
			}
			noteInfo.PinnedDate = pinnedDate
		}
		if err == nil && res == 1 {
			noteInfo.Pinned = true
			if NoteWidgets.pinButton != nil {
				NoteWidgets.pinButton.SetIcon(theme.RadioButtonCheckedIcon())
				NoteWidgets.pinButton.Refresh()
			}
		}
	}

	if AppStatus.currentView == VIEW_PINNED {
		UpdateView() //updates view on main window
	}
	UpdateProperties(noteInfo)
}

func SetEditMode() {
	NoteWidgets.markdownText.Hide()
	NoteContainers.markdown.Hide()
	NoteWidgets.deleteButton.Show()
	NoteWidgets.modeSelect.SetSelected(EDIT_MODE)
	NoteWidgets.entry.Show()
	noteWindow.Canvas().Focus(NoteWidgets.entry) //this seems to make no difference!!!
	//noteWindow.Content().Refresh()
}

func SetViewMode() {
	NoteWidgets.entry.Hide()
	NoteWidgets.deleteButton.Hide()
	NoteWidgets.markdownText.ParseMarkdown(NoteWidgets.entry.Text)
	NoteWidgets.markdownText.Show()
	NoteWidgets.modeSelect.SetSelected(VIEW_MODE)
	noteWindow.Canvas().Focus(nil) // this allows the canvas keyboard shortcuts to work rather than the entry widget shortcuts
	NoteContainers.markdown.Show()
}

func ChangeNoteColour(noteInfo *note.NoteInfo) {
	picker := dialog.NewColorPicker("Note Color", "Pick colour", func(c color.Color) {
		fmt.Println(c)
		hex := FyneColourToRGBHex(c)
		noteInfo.Colour = fmt.Sprintf("%s%s", "#", hex)
		//noteColour, err := RGBStringToFyneColor(fmt.Sprintf("%s%s", "#", hex))
		//if err != nil {
		//	log.Panicln(err)
		//}
		NoteCanvas.noteBackground.FillColor = c
	}, noteWindow)
	picker.Advanced = true
	picker.Show()
	UpdateProperties(noteInfo)
}

func AddNoteKeyboardShortcuts(noteInfo *note.NoteInfo) {
	//Keyboard shortcut to set edit mode
	ctrl_e := &desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierControl,
	}

	noteWindow.Canvas().AddShortcut(ctrl_e, func(shortcut fyne.Shortcut) {
		SetEditMode()
	})

	//Keyboard shortcut to set view mode
	ctrl_q := &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}

	noteWindow.Canvas().AddShortcut(ctrl_q, func(shortcut fyne.Shortcut) {
		SetViewMode()
	})

	//Keyboard shortcut to pin/unpin notes
	ctrl_p := &desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl,
	}

	noteWindow.Canvas().AddShortcut(ctrl_p, func(shortcut fyne.Shortcut) {
		PinNote(noteInfo)
	})

	//Keyboard shortcut to change note colour
	ctrl_h := &desktop.CustomShortcut{
		KeyName:  fyne.KeyH,
		Modifier: fyne.KeyModifierControl,
	}

	noteWindow.Canvas().AddShortcut(ctrl_h, func(shortcut fyne.Shortcut) {
		ChangeNoteColour(noteInfo)
	})
}

func ShowProperties(noteInfo *note.NoteInfo) {
	if NoteContainers.propertiesPanel.Hidden {
		text := note.GetPropertiesText(noteInfo)
		NoteWidgets.propertiesText.SetText(text)
		NoteContainers.propertiesPanel.Show()
	} else {
		NoteContainers.propertiesPanel.Hide()
	}
}

func UpdateProperties(noteInfo *note.NoteInfo) {
	if !NoteContainers.propertiesPanel.Hidden {
		text := note.GetPropertiesText(noteInfo)
		NoteWidgets.propertiesText.SetText(text)
		NoteWidgets.propertiesText.Refresh()
	}
}
