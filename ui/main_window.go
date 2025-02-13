package ui

import (
	//"image/color"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	 _"fmt"
)

var listPanel *fyne.Container
var grid *fyne.Container
var noteSize fyne.Size = fyne.NewSize(500,400) //default note size, may be overidden by user prefs
var recentNotesLimit = 6 //default,  may be overidden by user prefs


func StartUI(){

	app := app.New()
	CreateMainWindow(app)
}

func CreateMainWindow(app fyne.App){

	mainWindow := app.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	noteSize = fyne.NewSize(500,400) //this should depend on resolution of current display
	recentNotesLimit = 6
	//initialView := "Recent"
	//**************************************************************

	//Main Grid container for displaying notes
	//grid = container.New(layout.NewGridWrapLayout(noteSize))

	grid = container.NewGridWrap(noteSize)

	//Create the side panel
	//side := CreateSidePanel()
	//SIDE PANEL-------------------------------------------------


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
		notes,_ := scribedb.GetPinnedNotes()
		ShowNotesInGrid(grid,notes,noteSize)
	})

	RecentBtn := widget.NewButton("R", func(){
		if listPanel != nil{
			listPanel.Hide()
		}
		notes,_ := scribedb.GetRecentNotes(recentNotesLimit)
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

	notebooks,_ := scribedb.GetNotebooks()

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
		richText := widget.NewRichTextFromMarkdown(note.Content)
		richText.Wrapping = fyne.TextWrapWord
		bgColour, _ := RGBStringToFyneColor(note.BackgroundColour)
		noteColourRect := canvas.NewRectangle(bgColour) // this is the not colour marker (used to be note background in old scribe)
		colourLabel := canvas.NewText("", bgColour) // this is only used to size the note colour rectangle
		colourLabel.TextSize = 13
		contStacked := container.NewStack(colourLabel,noteColourRect) //stacked sowe can use a coloured rectangle as the background to the label
		cont := container.NewVBox(contStacked, richText)
		cont.Resize(noteSize)
		srcont := container.NewHScroll(cont)
		srcont.Resize(noteSize)
		grid.Add(srcont)
	}
}
