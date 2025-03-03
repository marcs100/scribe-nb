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

	//"fyne.io/fyne/v2/layout"
	"errors"
	"fmt"
	"log"
	"slices"
)

//var listPanel *fyne.Container
var noteWidth float32 = 500
var noteHeight float32 = 350
var noteSize fyne.Size = fyne.NewSize(noteWidth,noteHeight) //default note size, may be overidden by user prefs
//var toolbar *widget.Toolbar
var mainApp fyne.App
var viewLabel *widget.Label

//var noteBorderSize = fyne.NewSize(noteWidth+3,noteHeight+3)
var recentNotesLimit = 6 //default,  may be overidden by user prefs

var openNotes []uint //maintain a list of notes that are currently open

const VIEW_PINNED string = "pinned"
const VIEW_RECENT string = "recent"
const VIEW_NOTEBOOK string = "notebooks"
const VIEW_TAGS string = "tag"
const LAYOUT_GRID = "grid"
const LAYOUT_PAGE = "page"


func StartUI(){

	mainApp = app.New()
	CreateMainWindow()
}

func CreateMainWindow(){

	modDarkColour,_ := RGBStringToFyneColor("#2f2f2f")
	modLightColour,_ := RGBStringToFyneColor("#e2e2e2")

	themeVariant := mainApp.Settings().ThemeVariant()
	themeBgColour := mainApp.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)
	appStatus.themeBgColour = themeBgColour
	if themeVariant == theme.VariantDark{
		appStatus.themeBgColour = modDarkColour
	}else if themeVariant == theme.VariantLight{
		appStatus.themeBgColour = modLightColour
	}

	mainWindow := mainApp.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	//noteSize = fyne.NewSize(500,400) //this should depend on resolution of current display
	recentNotesLimit = 6
	initialView := VIEW_PINNED
	initialLayout := LAYOUT_GRID
	//**************************************************************

	pageView.CurrentPage = 0
	pageView.NumberOfPages = 0

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(noteSize)
	appContainers.grid = grid //store to allow interaction in otehr functions

	singleNotePage := widget.NewRichTextFromMarkdown("")
	appWidgets.singleNotePage = singleNotePage
	singleNoteStack := container.NewStack()
	appContainers.singleNoteStack = singleNoteStack

	//Create The main panel
	main := CreateMainPanel()

	top := CreateTopPanel()

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(top, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))

	//set default view and layout`
	appStatus.currentView = initialView
	appStatus.currentLayout = initialLayout

	UpdateView()

	mainWindow.ShowAndRun()
}


func CreateMainPanel()(*fyne.Container){

	mainGridContainer := container.NewScroll(appContainers.grid)
	appContainers.mainGridContainer = mainGridContainer
	mainPageContainer := container.NewScroll(appContainers.singleNoteStack)
	appContainers.mainPageContainer = mainPageContainer
	mainStackedContainer := container.NewStack(mainPageContainer,mainGridContainer )

	return mainStackedContainer
}

func CreateTopPanel()(*fyne.Container){
	viewLabel = widget.NewLabel("Pinned Notes")
	viewLabelFixed := widget.NewLabel("Viewing: ")
	spacer := widget.NewLabel("    ")



	toolbar  := widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func(){
			if appStatus.currentLayout != LAYOUT_GRID{
				appStatus.currentLayout = LAYOUT_GRID
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.DocumentIcon(), func(){
			if appStatus.currentLayout != LAYOUT_PAGE{
				appStatus.currentLayout = LAYOUT_PAGE
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.NavigateBackIcon(), func(){
			if appStatus.currentLayout == LAYOUT_PAGE{
				if pageView.PageBack() > 0{
					UpdateView()
				}
			}
		}),
		widget.NewToolbarAction(theme.NavigateNextIcon(), func(){
			if appStatus.currentLayout == LAYOUT_PAGE{
				if pageView.PageForward() > 0{
					UpdateView()
				}

			}
		}),
	)

	appWidgets.toolbar = toolbar
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
		appStatus.notes,err = scribedb.GetPinnedNotes()
		if err != nil{
			log.Print("Error getting pinned notes: ")
			log.Panic(err)
		}
		appStatus.currentView = VIEW_PINNED
		//ShowNotesInGrid(notes,noteSize)
		pageView.Reset()
		UpdateView()
	})

	RecentBtn := widget.NewButton("R", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		var err error
		appStatus.notes,err = scribedb.GetRecentNotes(recentNotesLimit)
		if err != nil{
			log.Print("Error getting recent notes: ")
			log.Panic(err)
		}
		appStatus.currentView = VIEW_RECENT
		//ShowNotesInGrid(notes,noteSize)
		pageView.Reset()
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

			if appContainers.grid != nil{
				appContainers.grid.RemoveAll()
			}
		}
		pageView.Reset()
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
				appStatus.notes,_ = scribedb.GetNotebook(notebooks[id])
				appStatus.currentView = VIEW_NOTEBOOK
				appStatus.currentNotebook = notebooks[id]
				//ShowNotesInGrid(notes, noteSize)
				pageView.Reset()
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
	if appContainers.grid == nil || appContainers.mainGridContainer == nil{
		return
	}

	if appContainers.mainPageContainer != nil{
		appContainers.mainPageContainer.Hide()
	}

	appContainers.grid.RemoveAll()
	for _, note := range notes{
		richText := newScribeNoteText(note.Content, func(){
			//fmt.Println("You clciked note with id ... " + fmt.Sprintf("%d", note.Id))
			if slices.Contains(openNotes, note.Id){
				//note is already open
				fmt.Println("note is already open")
			}else{
				openNotes = append(openNotes, note.Id)
				OpenNoteWindow(note.Id)
			}
		})
		richText.Wrapping = fyne.TextWrapWord
		themeBackground := canvas.NewRectangle(appStatus.themeBgColour)
		noteColour,_ := RGBStringToFyneColor(note.BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if note.BackgroundColour == "#e7edef" || note.BackgroundColour == "#FFFFFF"{
			noteBackground = canvas.NewRectangle(appStatus.themeBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack:= container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		appContainers.grid.Add(noteStack)
	}
	appContainers.grid.Refresh()
	appContainers.mainGridContainer.Show()
}

func ShowNotesAsPages(notes []scribedb.NoteData){
	if appContainers.mainGridContainer != nil{
		appContainers.mainGridContainer.Hide()
	}

	pageView.NumberOfPages = len(notes)
	if pageView.CurrentPage ==0{
		pageView.CurrentPage = 1
	}

	retreievdNote, err := scribedb.GetNote(notes[pageView.CurrentPage-1].Id)

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

	appWidgets.singleNotePage.ParseMarkdown(noteInfo.Content)
	appWidgets.singleNotePage.Wrapping = fyne.TextWrapWord
	appWidgets.singleNotePage.Refresh()

	themeBackground := canvas.NewRectangle(appStatus.themeBgColour)
	noteColour,_ := RGBStringToFyneColor(noteInfo.Colour)
	noteBackground := canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF"{
		noteBackground = canvas.NewRectangle(appStatus.themeBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(noteBackground)
	textPadded := container.NewPadded(themeBackground, appWidgets.singleNotePage)
	noteStack:= container.NewStack(colourStack, textPadded)

	appContainers.singleNoteStack.RemoveAll()
	appContainers.singleNoteStack.Add(noteStack)

	appContainers.mainPageContainer.Show()
	appContainers.mainPageContainer.Refresh()

}


func UpdateView()error{
	//var notes []scribedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch appStatus.currentView{
		case VIEW_PINNED:
			viewLabel.SetText("Pinned Notes")
			appStatus.notes, err = scribedb.GetPinnedNotes()
		case VIEW_RECENT:
			viewLabel.SetText(("Recent Notes"))
			appStatus.notes, err = scribedb.GetRecentNotes(recentNotesLimit)
		case VIEW_NOTEBOOK:
			viewLabel.SetText("Notebook - " + appStatus.currentNotebook)
			appStatus.notes, err = scribedb.GetNotebook(appStatus.currentNotebook)
		default:
			err = errors.New("undefined view")
	}

	if err != nil{
		return err
	}

	switch appStatus.currentLayout{
		case LAYOUT_GRID:
			appWidgets.toolbar.Items[2].ToolbarObject().Hide()
			appWidgets.toolbar.Items[3].ToolbarObject().Hide()
			ShowNotesInGrid(appStatus.notes, noteSize)
		case LAYOUT_PAGE:
			appWidgets.toolbar.Items[2].ToolbarObject().Show()
			appWidgets.toolbar.Items[3].ToolbarObject().Show()
			ShowNotesAsPages(appStatus.notes)
	}

	return err
}
