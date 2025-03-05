package ui

import (
	//"image/color"
	//"strconv"
	"scribe-nb/scribedb"
	"scribe-nb/note"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	//"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"scribe-nb/config"
	//"fyne.io/fyne/v2/layout"
	"errors"
	"fmt"
	"log"
	"slices"
)

//var listPanel *fyne.Container
//var noteWidth float32 = 500
//var noteHeight float32 = 350
//var toolbar *widget.Toolbar
var mainApp fyne.App
var viewLabel *widget.Label

//var noteBorderSize = fyne.NewSize(noteWidth+3,noteHeight+3)
//var recentNotesLimit = 6 //default,  may be overidden by user prefs

const VIEW_PINNED string = "pinned"
const VIEW_RECENT string = "recent"
const VIEW_NOTEBOOK string = "notebooks"
const VIEW_TAGS string = "tag"
const LAYOUT_GRID = "grid"
const LAYOUT_PAGE = "page"


func StartUI(appConfigIn *config.Config){

	Conf = appConfigIn
	mainApp = app.New()
	CreateMainWindow()
}

func CreateMainWindow(){

	modDarkColour,_ := RGBStringToFyneColor("#2f2f2f")
	modLightColour,_ := RGBStringToFyneColor("#e2e2e2")

	themeVariant := mainApp.Settings().ThemeVariant()
	themeBgColour := mainApp.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)
	AppStatus.themeBgColour = themeBgColour
	if themeVariant == theme.VariantDark{
		AppStatus.themeBgColour = modDarkColour
	}else if themeVariant == theme.VariantLight{
		AppStatus.themeBgColour = modLightColour
	}


	AppStatus.noteSize = fyne.NewSize(Conf.AppSettings.NoteWidth,Conf.AppSettings.NoteHeight)

	mainWindow := mainApp.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	//noteSize = fyne.NewSize(500,400) //this should depend on resolution of current display
	//recentNotesLimit = 6
	//initialView := VIEW_PINNED
	//initialLayout := LAYOUT_GRID
	//**************************************************************

	PageView.CurrentPage = 0
	PageView.NumberOfPages = 0

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(AppStatus.noteSize)
	AppContainers.grid = grid //store to allow interaction in otehr functions

	singleNotePage := widget.NewRichTextFromMarkdown("")
	AppWidgets.singleNotePage = singleNotePage
	singleNoteStack := container.NewStack()
	AppContainers.singleNoteStack = singleNoteStack

	//Create The main panel
	main := CreateMainPanel()

	top := CreateTopPanel()

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(top, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))

	//set default view and layout`
	AppStatus.currentView = Conf.AppSettings.InitialView
	AppStatus.currentLayout = Conf.AppSettings.InitialLayout

	UpdateView()

	mainWindow.ShowAndRun()
}


func CreateMainPanel()(*fyne.Container){

	mainGridContainer := container.NewScroll(AppContainers.grid)
	AppContainers.mainGridContainer = mainGridContainer
	mainPageContainer := container.NewScroll(AppContainers.singleNoteStack)
	AppContainers.mainPageContainer = mainPageContainer
	mainStackedContainer := container.NewStack(mainPageContainer,mainGridContainer )

	return mainStackedContainer
}

func CreateTopPanel()(*fyne.Container){
	viewLabel = widget.NewLabel("Pinned Notes")
	viewLabelFixed := widget.NewLabel("Viewing: ")
	spacer := widget.NewLabel("    ")



	toolbar  := widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func(){
			if AppStatus.currentLayout != LAYOUT_GRID{
				AppStatus.currentLayout = LAYOUT_GRID
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.DocumentIcon(), func(){
			if AppStatus.currentLayout != LAYOUT_PAGE{
				AppStatus.currentLayout = LAYOUT_PAGE
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.NavigateBackIcon(), func(){
			if AppStatus.currentLayout == LAYOUT_PAGE{
				if PageView.PageBack() > 0{
					UpdateView()
				}
			}
		}),
		widget.NewToolbarAction(theme.NavigateNextIcon(), func(){
			if AppStatus.currentLayout == LAYOUT_PAGE{
				if PageView.PageForward() > 0{
					UpdateView()
				}

			}
		}),
	)

	AppWidgets.toolbar = toolbar
	topPanel := container.NewHBox(spacer, viewLabelFixed, viewLabel, toolbar)
	return topPanel
}

func CreateSidePanel()(*fyne.Container){
	var listPanel *fyne.Container
	pinnedBtn := widget.NewButton("P", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		var err error
		AppStatus.notes,err = scribedb.GetPinnedNotes()
		if err != nil{
			log.Print("Error getting pinned notes: ")
			log.Panic(err)
		}
		AppStatus.currentView = VIEW_PINNED
		//ShowNotesInGrid(notes,noteSize)
		PageView.Reset()
		UpdateView()
	})

	RecentBtn := widget.NewButton("R", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		var err error
		AppStatus.notes,err = scribedb.GetRecentNotes(Conf.AppSettings.RecentNotesLimit)
		if err != nil{
			log.Print("Error getting recent notes: ")
			log.Panic(err)
		}
		AppStatus.currentView = VIEW_RECENT
		//ShowNotesInGrid(notes,noteSize)
		PageView.Reset()
		UpdateView()
	})

	notebooksBtn := widget.NewButton("N", func(){
		viewLabel.SetText("Notebooks")
		if listPanel != nil{
			if listPanel.Visible(){
				listPanel.Hide()
			}else{
				listPanel.Show()
			}

			if AppContainers.grid != nil{
				AppContainers.grid.RemoveAll()
			}
		}
		PageView.Reset()
	})


	notebooks,err := scribedb.GetNotebooks()
	if err != nil{
		log.Print("Error getting Notebooks: ")
		log.Panic(err)
	}
	/*nbCovers,err := scribedb.GetNotebookCovers()
	if err != nil{
		log.Print("Error getting notrbook covers: ")
		log.Panic(err)
	}
	sort.Strings(nbCovers)*/

	notebooksList := widget.NewList(
		func()int {
			return len(notebooks)
		},
		func() fyne.CanvasObject{
			return widget.NewButton("------------Notebooks------------", func(){})

		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(notebooks[id])
			o.(*widget.Button).OnTapped = func(){
				AppStatus.notes,_ = scribedb.GetNotebook(notebooks[id])
				AppStatus.currentView = VIEW_NOTEBOOK
				AppStatus.currentNotebook = notebooks[id]
				//ShowNotesInGrid(notes, noteSize)
				PageView.Reset()
				UpdateView()
			}
		},
	)

	btnPanel := container.NewVBox(pinnedBtn, RecentBtn, notebooksBtn)
	listPanel = container.NewStack(notebooksList)
	listPanel.Hide()


	sideContainer := container.NewHBox(btnPanel,listPanel)

	return sideContainer
}


func ShowNotesInGrid(notes []scribedb.NoteData, noteSize fyne.Size){
	if AppContainers.grid == nil || AppContainers.mainGridContainer == nil{
		return
	}

	if AppContainers.mainPageContainer != nil{
		AppContainers.mainPageContainer.Hide()
	}

	AppContainers.grid.RemoveAll()
	for _, note := range notes{
		richText := newScribeNoteText(note.Content, func(){
			//fmt.Println("You clciked note with id ... " + fmt.Sprintf("%d", note.Id))
			if slices.Contains(AppStatus.openNotes, note.Id){
				//note is already open
				fmt.Println("note is already open")
			}else{
				AppStatus.openNotes = append(AppStatus.openNotes, note.Id)
				OpenNoteWindow(note.Id)
			}
		})
		richText.Wrapping = fyne.TextWrapWord
		themeBackground := canvas.NewRectangle(AppStatus.themeBgColour)
		noteColour,_ := RGBStringToFyneColor(note.BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if note.BackgroundColour == "#e7edef" || note.BackgroundColour == "#FFFFFF"{
			noteBackground = canvas.NewRectangle(AppStatus.themeBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack:= container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		AppContainers.grid.Add(noteStack)
	}
	AppContainers.grid.Refresh()
	AppContainers.mainGridContainer.Show()
}

func ShowNotesAsPages(notes []scribedb.NoteData){
	if AppContainers.mainGridContainer != nil{
		AppContainers.mainGridContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	if PageView.CurrentPage ==0{
		PageView.CurrentPage = 1
	}

	retreievdNote, err := scribedb.GetNote(notes[PageView.CurrentPage-1].Id)

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

	/*if noteInfo.Id != 0{
		noteInfo.NewNote = false
	}else{
		noteInfo.NewNote = true
	}

	if retreievdNote.Pinned > 0{
		noteInfo.Pinned = true
	} else {
		noteInfo.Pinned = false
	}*/

	//calculate initial note content hash
	//note.UpdateHash(&noteInfo)

	AppWidgets.singleNotePage.ParseMarkdown(noteInfo.Content)
	AppWidgets.singleNotePage.Wrapping = fyne.TextWrapWord
	AppWidgets.singleNotePage.Refresh()

	themeBackground := canvas.NewRectangle(AppStatus.themeBgColour)
	noteColour,_ := RGBStringToFyneColor(noteInfo.Colour)
	noteBackground := canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF"{
		noteBackground = canvas.NewRectangle(AppStatus.themeBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(noteBackground)
	textPadded := container.NewPadded(themeBackground, AppWidgets.singleNotePage)
	noteStack:= container.NewStack(colourStack, textPadded)

	AppContainers.singleNoteStack.RemoveAll()
	AppContainers.singleNoteStack.Add(noteStack)

	AppContainers.mainPageContainer.Show()
	AppContainers.mainPageContainer.Refresh()

}


func UpdateView()error{
	//var notes []scribedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch AppStatus.currentView{
		case VIEW_PINNED:
			viewLabel.SetText("Pinned Notes")
			AppStatus.notes, err = scribedb.GetPinnedNotes()
		case VIEW_RECENT:
			viewLabel.SetText(("Recent Notes"))
			AppStatus.notes, err = scribedb.GetRecentNotes(Conf.AppSettings.RecentNotesLimit)
		case VIEW_NOTEBOOK:
			viewLabel.SetText("Notebook - " + AppStatus.currentNotebook)
			AppStatus.notes, err = scribedb.GetNotebook(AppStatus.currentNotebook)
		default:
			err = errors.New("undefined view")
	}

	if err != nil{
		return err
	}

	switch AppStatus.currentLayout{
		case LAYOUT_GRID:
			AppWidgets.toolbar.Items[2].ToolbarObject().Hide()
			AppWidgets.toolbar.Items[3].ToolbarObject().Hide()
			ShowNotesInGrid(AppStatus.notes, AppStatus.noteSize)
		case LAYOUT_PAGE:
			AppWidgets.toolbar.Items[2].ToolbarObject().Show()
			AppWidgets.toolbar.Items[3].ToolbarObject().Show()
			ShowNotesAsPages(AppStatus.notes)
	}

	return err
}
