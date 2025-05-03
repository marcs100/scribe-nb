package ui

import (
	"fmt"
	"image/color"
	"log"
	"scribe-nb/note"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewNoteContainer(noteId uint, noteInfo *note.NoteInfo, retrievedNote *scribedb.NoteData, allowEdit bool, parentWindow fyne.Window) *fyne.Container {
	//var err error
	//var retrievedNote scribedb.NoteData
	//var noteInfo note.NoteInfo

	if noteId == 0 {
		//New note
		*noteInfo = note.NoteInfo{
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
		/*=retrievedNote, err = scribedb.GetNote(noteId)

		if err != nil {
			log.Println("error getting note")
			dialog.ShowError(err, noteWindow)
			log.Panic(err)
			}*/

		*noteInfo = note.NoteInfo{
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
	note.UpdateHash(noteInfo)

	//setup keyboard shortcuts
	NoteWidgets.entry = NewEntryCustom(func(cs *desktop.CustomShortcut) {
		switch cs.ShortcutName() {
		case ctrl_shift_q.ShortcutName():
			SetViewMode(parentWindow)
		case ctrl_shift_p.ShortcutName():
			PinNote(noteInfo)
		case ctrl_shift_c.ShortcutName():
			ChangeNoteColour(noteInfo, parentWindow)
		case ctrl_shift_i.ShortcutName():
			go ShowProperties(noteInfo)
		}
	}, func() { SaveNote(noteInfo, retrievedNote) })

	//NoteWidgets.entry = widget.NewMultiLineEntry()
	NoteWidgets.entry.Text = noteInfo.Content
	NoteWidgets.entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
	noteColour, _ := RGBStringToFyneColor(noteInfo.Colour)

	NoteCanvas.noteBackground = canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF" || noteInfo.Colour == "#ffffff" || noteInfo.Colour == "#000000" {
		NoteCanvas.noteBackground = canvas.NewRectangle(AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(NoteCanvas.noteBackground)

	NoteWidgets.markdownText = widget.NewRichTextFromMarkdown(noteInfo.Content)
	/*NoteWidgets.markdownText = NewRichTextFromMarkdownCustom(noteInfo.Content, func(cs *desktop.CustomShortcut) {
	switch cs.ShortcutName() {
	case ctrl_e.ShortcutName():
		SetEditMode(parentWindow)
	case ctrl_p.ShortcutName():
		PinNote(noteInfo)
	case ctrl_h.ShortcutName():
		ChangeNoteColour(noteInfo, parentWindow)
	}
	})*/
	NoteWidgets.markdownText.Wrapping = fyne.TextWrapWord
	NoteWidgets.markdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, NoteWidgets.markdownText)
	NoteContainers.markdown = container.NewStack(colourStack, markdownPadded)
	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(NoteContainers.markdown)
	background := canvas.NewRectangle(AppTheme.NoteBgColour)
	content := container.NewStack(background, scrolledMarkdown, NoteWidgets.entry)

	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if noteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	NoteWidgets.pinButton = widget.NewButtonWithIcon("", btnIcon, func() {
		PinNote(noteInfo)
	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	changeNotebookBtn := NewChangeNotebookButton(noteInfo, parentWindow)

	colourButton := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		ChangeNoteColour(noteInfo, parentWindow)
	})

	NoteWidgets.deleteButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		DeleteNote(noteInfo, parentWindow)
	})

	propertiesButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() { go ShowProperties(noteInfo) })

	NoteWidgets.deleteButton.Hide()

	NoteWidgets.modeSelect = widget.NewRadioGroup([]string{EDIT_MODE, VIEW_MODE}, func(value string) {
		switch value {
		case EDIT_MODE:
			if allowEdit {
				SetEditMode(parentWindow)
			}
		case VIEW_MODE:
			SetViewMode(parentWindow)
		}
	})

	if !allowEdit {
		NoteWidgets.modeSelect.Hide()
	}

	NoteContainers.propertiesPanel = NewProperetiesPanel()

	NoteWidgets.modeSelect.SetSelected("View")
	NoteWidgets.modeSelect.Horizontal = true
	toolbar := container.NewHBox(NoteWidgets.modeSelect, spacerLabel, NoteWidgets.pinButton, colourButton, changeNotebookBtn, propertiesButton, NoteWidgets.deleteButton)
	NoteContainers.propertiesPanel.Hide()

	return container.NewBorder(toolbar, nil, nil, NoteContainers.propertiesPanel, content)
}

func NewChangeNotebookButton(noteInfo *note.NoteInfo, parentWindow fyne.Window) *widget.Button {
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
							parentWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
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
			}, parentWindow)
			newNotebookDialog.Show()
		})

		nbMenu.Items = append(nbMenu.Items, nbMenuItem)

		//Now add all the existing notebooks to the menu
		for _, notebook := range notebooks {
			menuItem := fyne.NewMenuItem(notebook, func() {
				noteInfo.Notebook = notebook
				//fmt.Println("Change notebook to " + notebook)
				parentWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
				UpdateProperties(noteInfo)
			})
			nbMenu.Items = append(nbMenu.Items, menuItem)
		}

		popUpMenu := widget.NewPopUpMenu(nbMenu, parentWindow.Canvas())
		//popUpMenu.Show()
		pos := fyne.NewPos(250, 40)
		popUpMenu.ShowAtPosition(pos)
		//popUpMenu.ShowAtPosition(e.Position.AddXY(150,0))

	})

	return changeNotebookBtn
}

func DeleteNote(noteInfo *note.NoteInfo, parentWindow fyne.Window) {
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
				parentWindow.Close()
			}
		}
	}, parentWindow)
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
				log.Println(fmt.Sprintf("Error getting pinned date: %s", err))
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

	UpdateProperties(noteInfo)

	if AppStatus.currentView == VIEW_PINNED {
		UpdateView() //updates view on main window
	}
}

func SetEditMode(parentWindow fyne.Window) {
	NoteWidgets.markdownText.Hide()
	NoteContainers.markdown.Hide()
	NoteWidgets.deleteButton.Show()
	if AppStatus.currentLayout == LAYOUT_PAGE {
		//Hide page back & forward for edit mode
		AppWidgets.toolbar.Items[2].ToolbarObject().Hide()
		AppWidgets.toolbar.Items[3].ToolbarObject().Hide()
	}
	NoteWidgets.modeSelect.SetSelected(EDIT_MODE)
	NoteWidgets.entry.Show()
	parentWindow.Canvas().Focus(NoteWidgets.entry)
	parentWindow.Content().Refresh()
}

func SetViewMode(parentWindow fyne.Window) {
	NoteWidgets.entry.Hide()
	NoteWidgets.deleteButton.Hide()
	if AppStatus.currentLayout == LAYOUT_PAGE {
		//Show page back & forward for edit mode
		AppWidgets.toolbar.Items[2].ToolbarObject().Show()
		AppWidgets.toolbar.Items[3].ToolbarObject().Show()
	}
	NoteWidgets.markdownText.ParseMarkdown(NoteWidgets.entry.Text)
	NoteWidgets.markdownText.Show()
	NoteWidgets.modeSelect.SetSelected(VIEW_MODE)
	parentWindow.Canvas().Focus(nil) // this allows the canvas keyboard shortcuts to work rather than the entry widget shortcuts
	NoteContainers.markdown.Show()
}

func ChangeNoteColour(noteInfo *note.NoteInfo, parentWindow fyne.Window) {
	picker := dialog.NewColorPicker("Note Color", "Pick colour", func(c color.Color) {
		fmt.Println(c)
		hex := FyneColourToRGBHex(c)
		noteInfo.Colour = fmt.Sprintf("%s%s", "#", hex)
		//noteColour, err := RGBStringToFyneColor(fmt.Sprintf("%s%s", "#", hex))
		//if err != nil {
		//	log.Panicln(err)
		//}
		NoteCanvas.noteBackground.FillColor = c
	}, parentWindow)
	picker.Advanced = true
	picker.Show()
	UpdateProperties(noteInfo)
}

func SaveNote(noteInfo *note.NoteInfo, retrievedNote *scribedb.NoteData) {
	noteInfo.Content = NoteWidgets.entry.Text
	var noteChanges note.NoteChanges
	if noteInfo.NewNote {
		if noteInfo.Content != "" {
			noteChanges.ContentChanged = true
		}
	} else {
		noteChanges = note.CheckChanges(retrievedNote, noteInfo)
	}
	//if contentChanged{
	if noteChanges.ContentChanged || noteChanges.ParamsChanged {
		res, err := note.SaveNote(noteInfo)
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
		res, err := note.SaveNoteNoTimeStamp(noteInfo)
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
}

func AddNoteKeyboardShortcuts(noteInfo *note.NoteInfo, allowEdit bool, parentWindow fyne.Window) {
	//Keyboard shortcut to set edit mode
	if allowEdit {
		parentWindow.Canvas().AddShortcut(ctrl_shift_e, func(shortcut fyne.Shortcut) {
			SetEditMode(parentWindow)
		})
	}

	//Keyboard shortcut to pin/unpin notes
	parentWindow.Canvas().AddShortcut(ctrl_shift_p, func(shortcut fyne.Shortcut) {
		PinNote(noteInfo)
	})

	//Keyboard shortcut to change note colour
	parentWindow.Canvas().AddShortcut(ctrl_shift_c, func(shortcut fyne.Shortcut) {
		ChangeNoteColour(noteInfo, parentWindow)
	})

	//Keyboard shortcut to show properties panel
	parentWindow.Canvas().AddShortcut(ctrl_shift_i, func(shortcut fyne.Shortcut) {
		go ShowProperties(noteInfo)
	})
}
