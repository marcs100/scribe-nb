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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"fyne.io/fyne/v2/data/binding"
	//"github.com/fyne-io/terminal"
)

func OpenNoteWindow(noteId uint) {
	var PinBtn *widget.Button
	var err error
	var retrievedNote scribedb.NoteData
	var noteInfo note.NoteInfo

	if noteId == 0{
		//New note
		noteInfo = note.NoteInfo{
			Id: noteId,
			Notebook: "General",
			DateCreated: "",
			DateModified: "",
			Pinned: false,
			Colour: "#FFFFFF",
			Content: "",
			Deleted: false,
		}
	}else{
		//existing note
		retrievedNote, err = scribedb.GetNote(noteId)

		if err != nil{
			log.Println("error getting note")
			log.Panic(err)
		}

		noteInfo = note.NoteInfo{
			Id: retrievedNote.Id,
			Notebook: retrievedNote.Notebook,
			DateCreated: retrievedNote.Created,
			DateModified: retrievedNote.Modified,
			Colour: retrievedNote.BackgroundColour,
			Content: retrievedNote.Content,
			Deleted: false,
		}

		if retrievedNote.Pinned > 0{
			noteInfo.Pinned = true
		} else {
			noteInfo.Pinned = false
		}
	}


	if noteInfo.Id != 0{
		noteInfo.NewNote = false
	}else{
		noteInfo.NewNote = true
	}

	//calculate initial note content hash
	note.UpdateHash(&noteInfo)

	noteWindow := mainApp.NewWindow(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
	noteWindow.Resize(fyne.NewSize(900, 750))

	entry := widget.NewMultiLineEntry()
	entry.Text = noteInfo.Content
	entry.Wrapping = fyne.TextWrapWord

	themeBackground := canvas.NewRectangle(AppStatus.themeBgColour)
	noteColour,_ := RGBStringToFyneColor(noteInfo.Colour)

	noteBackground := canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF"{
		noteBackground = canvas.NewRectangle(AppStatus.themeBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(noteBackground)

	markdownText := widget.NewRichTextFromMarkdown(noteInfo.Content)
	markdownText.Wrapping = fyne.TextWrapWord
	markdownText.Hide()
	markdownPadded := container.NewPadded(themeBackground, markdownText)
	markdown:= container.NewStack(colourStack, markdownPadded)



	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(markdown)
	background := canvas.NewRectangle(AppStatus.themeBgColour)
	content := container.NewStack(background, scrolledMarkdown, entry)

	var win *fyne.Container


	//var btnLabel = "Pin"
	btnIcon := theme.RadioButtonIcon()
	if noteInfo.Pinned {
		btnIcon = theme.RadioButtonCheckedIcon()
		//btnLabel = "Unpin"
	}

	PinBtn = widget.NewButtonWithIcon("", btnIcon , func(){
		var res int64
		var err error = nil
		if noteInfo.Pinned{
			if noteInfo.NewNote{
				//new note that hasn't been saved yet'
				res = 1
			}else{
				res,err = scribedb.UnpinNote(noteInfo.Id)
			}

			if err == nil && res==1{
				noteInfo.Pinned = false
				if PinBtn != nil{
					PinBtn.SetIcon(theme.RadioButtonIcon())
					PinBtn.Refresh()
				}
			}
		}else{
			if noteInfo.Id == 0{
				//new note that hasn't been saved yet'
				res = 1
			}else{
				res,err = scribedb.PinNote(noteInfo.Id)
			}
			if err == nil && res==1{
				noteInfo.Pinned = true
				if PinBtn != nil{
					PinBtn.SetIcon(theme.RadioButtonCheckedIcon())
					PinBtn.Refresh()
				}
			}
		}

	})

	//changeNotebookBtn := NewButtonWithPos("Change Notebook", func(e *fyne.PointEvent){
	changeNotebookBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func(){
		var notebooks []string
		var err error
		if notebooks, err = scribedb.GetNotebooks(); err != nil{
			log.Println("Error getting notebook")
			log.Panicln(err)
		}
		nbMenu := fyne.NewMenu("Select Notebook")

		for _, notebook := range notebooks{
			menuItem := fyne.NewMenuItem(notebook, func(){
				noteInfo.Notebook = notebook
				//fmt.Println("Change notebook to " + notebook)
				noteWindow.SetTitle(fmt.Sprintf("Notebook: %s --- Note id: %d", noteInfo.Notebook, noteInfo.Id))
			})
			nbMenu.Items = append(nbMenu.Items, menuItem)
		}


		popUpMenu := widget.NewPopUpMenu(nbMenu, noteWindow.Canvas())
		//popUpMenu.Show()
		pos := fyne.NewPos(250,40)
		popUpMenu.ShowAtPosition(pos)
		//popUpMenu.ShowAtPosition(e.Position.AddXY(150,0))

	})

	colourButton := widget.NewButtonWithIcon("",theme.ColorPaletteIcon(), func(){
		picker := dialog.NewColorPicker("Note Color", "Pick colour", func(c color.Color){
			fmt.Println(c)
			hex := FyneColourToRGBHex(c)
			noteInfo.Colour = fmt.Sprintf("%s%s","#",hex)
			noteColour,err = RGBStringToFyneColor(fmt.Sprintf("%s%s","#",hex))
			if err != nil{
				log.Panicln(err)
			}
			noteBackground.FillColor = c
		} ,noteWindow)
		picker.Advanced = true
		picker.Show()
	})

	deleteBtn := widget.NewButtonWithIcon("",theme.DeleteIcon() , func(){
		dialog.ShowConfirm("Delete note","Are you sure?",func(confirm bool){
			if confirm{
				var res int64
				var err error = nil
				if noteInfo.NewNote{
					res = 1;
				}else{
					res, err = scribedb.DeleteNote(noteInfo.Id)
				}

				if res == 0 || err != nil{
					log.Println("Error deleing notes")
				}else{
					noteInfo.Deleted = true
					noteWindow.Close()
				}
			}
		}, noteWindow)
	})
	deleteBtn.Hide()

	modeWidget := widget.NewRadioGroup([]string{"Edit", "View"}, func(value string){
		switch value{
			case "Edit":
				markdownText.Hide()
				markdown.Hide()
				deleteBtn.Show()
				noteWindow.Canvas().Focus(entry) //this seems to make no difference!!!
				entry.Show()
			case "View":
				entry.Hide()
				deleteBtn.Hide()
				markdownText.ParseMarkdown(entry.Text)
				markdownText.Show()
				markdown.Show()
		}
	})

	modeWidget.SetSelected("View")
	modeWidget.Horizontal = true;
	toolbar := container.NewHBox(modeWidget,spacerLabel, PinBtn, colourButton, changeNotebookBtn, deleteBtn)
	win = container.NewBorder(toolbar, nil,nil,nil,content)

	noteWindow.SetContent(win)
	noteWindow.Canvas().Focus(entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", noteInfo.Id))
		noteInfo.Content = entry.Text
		var contentChanged, paramsChanged bool = false, false
		if noteInfo.NewNote{
			if noteInfo.Content != ""{
				contentChanged = true
			}
		}else{
			contentChanged, paramsChanged = note.CheckChanges(&retrievedNote,&noteInfo)
		}
		//if contentChanged{
		if contentChanged || paramsChanged{
			res, err := note.SaveNote(&noteInfo)
			if err != nil{
				log.Println("Error saving note")
				log.Panic()
			}

			if res == 0{
				log.Println("No note was saved (affected rows = 0)")
			}else{
				log.Println("....Note updates successfully....")
				go UpdateView()
			}
		}
		/*} else if paramsChanged{
			res, err := note.SaveNoteNoTimeStamp(&noteInfo)
			if err != nil{
				log.Println("Error saving note")
				log.Panic()
			}

			if res == 0{
				log.Println("No note was saved (affected rows = 0)")
			}else{
				log.Println("....Note updates successfully....")
				go UpdateView()
			}
		} */

		if index := slices.Index(AppStatus.openNotes,noteInfo.Id); index != -1{
			AppStatus.openNotes = slices.Delete(AppStatus.openNotes,index,index+1)
		}

	})
	noteWindow.Show()
}
