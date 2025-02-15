package ui

import (
	//"image/color"
	//"strconv"
	"fmt"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	_ "fmt"
	"log"

	//"sort"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var listPanel *fyne.Container
var grid *fyne.Container
var noteWidth float32 = 500
var noteHeight float32 = 350
var noteSize fyne.Size = fyne.NewSize(noteWidth,noteHeight) //default note size, may be overidden by user prefs
//var noteBorderSize = fyne.NewSize(noteWidth+3,noteHeight+3)
var recentNotesLimit = 6 //default,  may be overidden by user prefs


func StartUI(){

	app := app.New()
	CreateMainWindow(app)
}

func CreateMainWindow(app fyne.App){

	mainWindow := app.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	//noteSize = fyne.NewSize(500,400) //this should depend on resolution of current display
	recentNotesLimit = 6
	//initialView := "Recent"
	//**************************************************************

	//Main Grid container for displaying notes
	grid = container.NewGridWrap(noteSize)

	//Create The main panel
	main := CreateMainPanel(app, grid)

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(nil, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))


	mainWindow.ShowAndRun()
}


func CreateMainPanel(app fyne.App, grid *fyne.Container)(*fyne.Container){

	themeVariant := app.Settings().ThemeVariant()
	themeColour := app.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)
	modDarkColour,_ := RGBStringToFyneColor("#2f2f2f")
	modLightColour,_ := RGBStringToFyneColor("#e2e2e2")

	mainContainer := container.NewScroll(grid)
	mainBackground := canvas.NewRectangle(themeColour)
	if themeVariant == theme.VariantDark{
		mainBackground = canvas.NewRectangle(modDarkColour)
	} else if themeVariant == theme.VariantLight{
		mainBackground = canvas.NewRectangle(modLightColour)
	}
	mainStackedContainer := container.NewStack(mainBackground, mainContainer)

	return mainStackedContainer
}

/*func CreateTopPanel()(*fyne.Container){

	sidePanelBtn := widget.NewButton("panel")
}*/


func CreateSidePanel()(*fyne.Container){

	pinnedBtn := widget.NewButton("P", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		notes,err := scribedb.GetPinnedNotes()
		if err != nil{
			log.Print("Error getting pinned notes: ")
			log.Panic(err)
		}
		ShowNotesInGrid(grid,notes,noteSize)
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
		ShowNotesInGrid(grid,notes,noteSize)
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
				ShowNotesInGrid(grid, notes, noteSize)
			}
		},
	)

	btnPanel := container.NewVBox(pinnedBtn, RecentBtn, notebooksBtn)
	listPanel = container.NewStack(notebooksList)
	listPanel.Hide()


	sideContainer := container.NewHBox(btnPanel,listPanel)

	return sideContainer
}

func ShowNotesInGrid(grid *fyne.Container, notes []scribedb.NoteData, noteSize fyne.Size){
	if grid == nil{
		return
	}

	grid.RemoveAll()
	for _, note := range notes{
		richText := newScribeNoteText(note.Content, func(){
			fmt.Println("You clciked note with id ... " + fmt.Sprintf("%d", note.Id))
		})
		richText.Wrapping = fyne.TextWrapWord
		//bgColour, _ := RGBStringToFyneColor(note.BackgroundColour)
		//noteColourRect := canvas.NewRectangle(bgColour) // this is the note colour marker (used to be note background in old scribe)
		//noteColourRect.Resize(noteBorderSize)

		//textCont := container.NewStack(richText)
		//textCont.Resize(noteSize)
		grid.Add(richText)
	}
}
