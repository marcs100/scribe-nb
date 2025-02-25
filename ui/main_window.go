package ui

import (
	//"image/color"
	//"strconv"
	"image/color"
	"scribe-nb/scribedb"

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
var grid *fyne.Container
var mainApp fyne.App
//var noteBorderSize = fyne.NewSize(noteWidth+3,noteHeight+3)
var recentNotesLimit = 6 //default,  may be overidden by user prefs

var openNotes []uint //maintain a list of notes that are currently open

const VIEW_PINNED string = "pinned"
const VIEW_RECENT string = "recent"
const VIEW_NOTEBOOK string = "notebook"
const VIEW_TAGS string = "tag"
var currentView string = ""
var currentNotebook string = ""

var themeBgColour color.Color

func StartUI(){

	mainApp = app.New()
	CreateMainWindow()
}

func CreateMainWindow(){

	modDarkColour,_ := RGBStringToFyneColor("#2f2f2f")
	modLightColour,_ := RGBStringToFyneColor("#e2e2e2")

	themeVariant := mainApp.Settings().ThemeVariant()
	themeBgColour = mainApp.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)
	if themeVariant == theme.VariantDark{
		themeBgColour = modDarkColour
	}else if themeVariant == theme.VariantLight{
		themeBgColour = modLightColour
	}

	mainWindow := mainApp.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	//noteSize = fyne.NewSize(500,400) //this should depend on resolution of current display
	recentNotesLimit = 6
	initialView := VIEW_PINNED
	//**************************************************************

	//Main Grid container for displaying notes
	grid = container.NewGridWrap(noteSize)

	//Create The main panel
	main := CreateMainPanel()

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(nil, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))

	//set defaukr view
	currentView = initialView
	UpdateView()

	mainWindow.ShowAndRun()
}


func CreateMainPanel()(*fyne.Container){

	mainContainer := container.NewScroll(grid)
	mainStackedContainer := container.NewStack(mainContainer)
	return mainStackedContainer
}

func CreateSidePanel()(*fyne.Container){
	var listPanel *fyne.Container
	pinnedBtn := widget.NewButton("P", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		notes,err := scribedb.GetPinnedNotes()
		if err != nil{
			log.Print("Error getting pinned notes: ")
			log.Panic(err)
		}
		currentView = VIEW_PINNED
		ShowNotesInGrid(notes,noteSize)
	})

	RecentBtn := widget.NewButton("R", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		notes,err := scribedb.GetRecentNotes(recentNotesLimit)
		if err != nil{
			log.Print("Error getting recent notes: ")
			log.Panic(err)
		}
		currentView = VIEW_RECENT
		ShowNotesInGrid(notes,noteSize)
	})

	notebooksBtn := widget.NewButton("N", func(){
		if listPanel != nil{
			if listPanel.Visible(){
				listPanel.Hide()
			}else{
				listPanel.Show()
			}

			if grid != nil{
				grid.RemoveAll()
			}
		}
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
				notes,_ := scribedb.GetNotebook(notebooks[id])
				currentView = VIEW_NOTEBOOK
				currentNotebook = notebooks[id]
				ShowNotesInGrid(notes, noteSize)
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
	if grid == nil{
		return
	}

	grid.RemoveAll()
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
		themeBackground := canvas.NewRectangle(themeBgColour)
		noteColour,_ := RGBStringToFyneColor(note.BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if note.BackgroundColour == "#e7edef" || note.BackgroundColour == "#FFFFFF"{
			noteBackground = canvas.NewRectangle(themeBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack:= container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		grid.Add(noteStack)
	}
	grid.Refresh()
}



func UpdateView()error{
	var notes []scribedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch currentView{
		case VIEW_PINNED:
			notes, err = scribedb.GetPinnedNotes()
		case VIEW_RECENT:
			notes, err = scribedb.GetRecentNotes(recentNotesLimit)
		case VIEW_NOTEBOOK:
			notes, err = scribedb.GetNotebook(currentNotebook)
		default:
			err = errors.New("undefined view")
	}

	if err == nil{
		ShowNotesInGrid(notes, noteSize)
	}

	return err
}
