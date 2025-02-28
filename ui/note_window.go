package ui

import (
	"fmt"
	"log"
	"scribe-nb/note"
	"scribe-nb/scribedb"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	//"github.com/fyne-io/terminal"
)

func OpenNoteWindow(noteId uint) {
	var PinBtn *widget.Button

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
		Deleted: false,
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

	noteWindow := mainApp.NewWindow(fmt.Sprintf("Notebook: %s --- Note id: %d", retreievdNote.Notebook, retreievdNote.Id))
	noteWindow.Resize(fyne.NewSize(900, 750))

	entry := widget.NewMultiLineEntry()
	entry.Text = noteInfo.Content
	entry.Wrapping = fyne.TextWrapWord

	markdown := widget.NewRichTextFromMarkdown(noteInfo.Content)
	markdown.Wrapping = fyne.TextWrapWord
	markdown.Hide()

	var btnLabel = "Pin"
	if noteInfo.Pinned {
		btnLabel = "Unpin"
	}

	fmt.Println("label is "+ btnLabel)

	PinBtn = widget.NewButton(btnLabel, func(){
		if noteInfo.Pinned{
			res,err := scribedb.UnpinNote(noteInfo.Id)
			if err == nil && res==1{
				noteInfo.Pinned = false
				if PinBtn != nil{
					PinBtn.SetText("Pin")
					PinBtn.Refresh()
				}
			}
		}else{
			res,err := scribedb.PinNote(noteInfo.Id)
			if err == nil && res==1{
				noteInfo.Pinned = true
				if PinBtn != nil{
					PinBtn.SetText("UnPin")
					PinBtn.Refresh()
				}
			}
		}

	})


	deleteBtn := widget.NewButton("Del", func(){
		dialog.ShowConfirm("Delete note","Are you sure?",func(confirm bool){
			if confirm{
				res, err := scribedb.DeleteNote(noteInfo.Id)
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
				markdown.Hide()
				deleteBtn.Show()
				noteWindow.Canvas().Focus(entry) //this seems to make no difference!!!
				entry.Show()
			case "View":
				entry.Hide()
				deleteBtn.Hide()
				markdown.ParseMarkdown(entry.Text)
				markdown.Show()
		}
	})

	modeWidget.SetSelected("View")
	modeWidget.Horizontal = true;

	spacerLabel := widget.NewLabel("      ")

	scrolledMarkdown := container.NewScroll(markdown)
	background := canvas.NewRectangle(themeBgColour)
	content := container.NewStack(background, entry, scrolledMarkdown)
	toolbar := container.NewHBox(modeWidget,spacerLabel, PinBtn, deleteBtn)
	win := container.NewBorder(toolbar, nil,nil,nil,content)

	noteWindow.SetContent(win)
	noteWindow.Canvas().Focus(entry)
	noteWindow.SetOnClosed(func() {
		fmt.Println(fmt.Sprintf("Closing note %d", noteInfo.Id))
		noteInfo.Content = entry.Text
		contentChanged, paramsChanged := note.CheckChanges(&retreievdNote,&noteInfo)
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

		if index := slices.Index(openNotes,noteInfo.Id); index != -1{
			openNotes = slices.Delete(openNotes,index,index+1)
		}

	})
	noteWindow.Show()
}
