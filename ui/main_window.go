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
	//"fmt"
)


func StartUI(){

	app := app.New()
	mainWindow := app.NewWindow("Scribe-NB")

	//Can we get the current display size??????????

	//******************************************************************
	//yellow := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	//red := color.NRGBA{R: 255,G: 0, B: 0, A: 255}
	//******************************************************************
	noteSize := fyne.NewSize(500,400) //this should depend on resolution of current display

	notes, _:= scribedb.GetPinnedNotes()

	//lets try adding a side pane
	/*
	"": {"A"},
	"A": {"B", "D"},
	"B": {"C"},
	"C": {"abc"},
	"D": {"E"},
	"E": {"F", "G"},
	*/
	nodes := map[string][]string{
		"": {"Views"},
		"Views": {"Pinned", "Recent", "Notebooks"},
	}

	tree := widget.NewTreeWithStrings(nodes)
	tree.OnSelected = func(id string) {
		println("Selected:", id)
	}
	tree.OnUnselected = func(id string) {
		println("Unselected:", id)
	}
	tree.OpenAllBranches()
	sideContainer := container.NewScroll(tree)
	sideContainer.SetMinSize(fyne.NewSize(200,200))


	//Display notes in a grid
	grid := container.New(layout.NewGridWrapLayout(noteSize))
	for _, note := range notes{
		richText := widget.NewRichTextFromMarkdown(note.Content)
		richText.Wrapping = fyne.TextWrapWord
		bgColour, _ := RGBStringToFyneColor(note.BackgroundColour)
		rect := canvas.NewRectangle(bgColour)
		//label := widget.NewLabel("")
		colourLabel := canvas.NewText("", bgColour)
		colourLabel.TextSize = 10
		contStacked := container.NewStack(colourLabel,rect)
		cont := container.NewVBox(contStacked, richText)
		cont.Resize(noteSize)
		srcont := container.NewHScroll(cont)
		srcont.Resize(noteSize)
		grid.Add(srcont)
	}

	mainContainer := container.NewScroll(grid)
	appContainer := container.NewBorder(nil, nil, sideContainer, nil, mainContainer)
	//appContainer := container.NewHSplit(sideContainer, mainContainer)
	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000,1200))

	tree.Refresh()
	tree.ScrollToBottom()

	mainWindow.ShowAndRun()
}
