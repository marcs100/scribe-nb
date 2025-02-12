package ui

import (
	//"image/color"
	"scribe-nb/scribedb"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	//"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fmt"
)

func StartUI(){

	app := app.New()
	CreateMainWindow(app)
}

func CreateMainWindow(app fyne.App){
	var grid *fyne.Container
	var viewsTree *widget.Tree
	var notebooksTree *widget.Tree
	var notebooks []string //maintain a list of notebook names
	//var tags []string //maintain a list of tags

	mainWindow := app.NewWindow("Scribe-NB")

	// Options that wil be part of config file ************************
	noteSize := fyne.NewSize(500,400) //this should depend on resolution of current display
	recentNotesLimit := 6
	initialView := "Recent"
	//**************************************************************

	//Main Grid container for displaying notes
	grid = container.New(layout.NewGridWrapLayout(noteSize))

	nodes := map[string][]string{
		"": {"Views"},
		"Views": {"Pinned", "Recent", "Notebooks"},
	}
	viewsTree = widget.NewTreeWithStrings(nodes)


	notebooks,_ = scribedb.GetNotebooks()
	nb_nodes := map[string][]string{
		"": {"notebooks"},
		"Notebooks": {},
	}

	nb_nodes["notebooks"] = notebooks
	notebooksTree = widget.NewTreeWithStrings(nb_nodes)

	//Create the side panel
	side := CreateSidePanel(app, viewsTree, notebooksTree, grid, noteSize, recentNotesLimit)

	//Create The main panel
	main := CreateMainPanel(app, grid)

	//layout the main window
	appContainer := container.NewBorder(nil, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))

	//Make the tree selection correspond to the initial view
	viewsTree.Select(initialView)

	mainWindow.ShowAndRun()
}


func CreateMainPanel(app fyne.App, grid *fyne.Container)(*fyne.Container){

	themeVariant := app.Settings().ThemeVariant()
	themeColour := app.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)
	modDarkColour,_ := RGBStringToFyneColor("#373737")
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


func CreateSidePanel(app fyne.App,
		     viewsTree *widget.Tree,
		     notebooksTree *widget.Tree,
		     grid *fyne.Container,
		     noteSize fyne.Size,
		     recentNotesLimit int)(*fyne.Container){

	themeVariant := app.Settings().ThemeVariant()
	themeColour := app.Settings().Theme().Color(theme.ColorNameBackground,themeVariant)

	modSideDarkColour,_ := RGBStringToFyneColor("#232323")
	modSideLightColour,_ := RGBStringToFyneColor("#f7e5e5")

	//lets try adding a side pane
	/*
	 " ": {"A"},                                                               *
	 "A": {"B", "D"},
	 "B": {"C"},
	 "C": {"abc"},
	 "D": {"E"},
	 "E": {"F", "G"},
	 */

	viewsTree.OnSelected = func(id string) {
		println("Selected:", id)
		switch id{
			case "Pinned":
				notes,_ := scribedb.GetPinnedNotes()
				ShowNotesInGrid(grid,notes, noteSize)
			case "Notebooks":
				notes,_ := scribedb.GetNotebook("General") //just for testing!!!!!!!
				ShowNotesInGrid(grid,notes, noteSize)
			case "Recent":
				notes,_ := scribedb.GetRecentNotes(recentNotesLimit)
				ShowNotesInGrid(grid,notes, noteSize)
		}
	}
	viewsTree.OnUnselected = func(id string) {
		println("Unselected:", id)
	}

	viewsTree.OpenAllBranches()

	scrollTreeView := container.NewScroll(viewsTree)
	scrollTreeView.SetMinSize(fyne.NewSize(200,200))

	scrollTreeNotebooks := container.NewVScroll(notebooksTree)
	scrollTreeNotebooks.SetMinSize(fyne.NewSize(200,500))

	//notebooksTree.OpenAllBranches()

	newNoteButton := widget.NewButton("New Note",func(){
		fmt.Println("New Note button pressed!")
	})

	spacerLabel := widget.NewLabel(" ")
	sideVBoxContainer := container.NewVBox(newNoteButton,spacerLabel, scrollTreeView, scrollTreeNotebooks)

	// set a background for the side pane depending on theme variant (light or dark)
	sideColourRect := canvas.NewRectangle(themeColour)
	if themeVariant == theme.VariantDark{
		sideColourRect = canvas.NewRectangle(modSideDarkColour)
	} else if themeVariant == theme.VariantLight{
		sideColourRect = canvas.NewRectangle(modSideLightColour)
	}

	sideStackedContainer := container.NewStack(sideColourRect,sideVBoxContainer)

	return sideStackedContainer
}

func ShowNotesInGrid(grid *fyne.Container, notes []scribedb.NoteData, noteSize fyne.Size){
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
