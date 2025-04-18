package ui

import (
	"errors"
	"fmt"
	"log"
	"scribe-nb/config"
	"scribe-nb/note"
	"scribe-nb/scribedb"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartUI(appConfigIn *config.Config, version string) {
	Conf = appConfigIn
	mainApp = app.NewWithID("scribe-nb")
	CreateMainWindow(version)
}

func CreateMainWindow(version string) {

	AppStatus.noteSize = fyne.NewSize(Conf.Settings.NoteWidth, Conf.Settings.NoteHeight)

	var themeVar theme_variant
	switch Conf.Settings.ThemeVariant {
	case "light":
		themeVar = LIGHT_THEME
	case "dark":
		themeVar = DARK_THEME
	case "system":
		themeVar = SYSTEM_THEME
	}

	AppTheme = GetThemeColours(themeVar)

	mainWindow = mainApp.NewWindow(fmt.Sprintf("Scribe-NB   v%s", version))

	//Main Grid container for displaying notes
	grid := container.NewGridWrap(AppStatus.noteSize)
	AppContainers.grid = grid //store to allow interaction in other functions

	singleNotePage := widget.NewRichTextFromMarkdown("")
	AppWidgets.singleNotePage = singleNotePage
	singleNoteStack := container.NewStack()
	AppContainers.singleNoteStack = singleNoteStack

	PageView.CurrentPage = 0
	PageView.NumberOfPages = 0

	//Create The main panel
	main := CreateMainPanel()

	top := CreateTopPanel()

	side := CreateSidePanel()

	//layout the main window
	appContainer := container.NewBorder(top, nil, side, nil, main)

	mainWindow.SetContent(appContainer)
	mainWindow.Resize(fyne.NewSize(2000, 1200))

	//set default view and layout`
	AppStatus.currentView = Conf.Settings.InitialView
	fmt.Println("initial view = " + Conf.Settings.InitialView)
	AppStatus.currentLayout = Conf.Settings.InitialLayout

	if err := UpdateView(); err != nil {
		fmt.Println(err)
	}

	//keyboard shortcuts
	AddMainKeyboardShortcuts()

	mainWindow.SetCloseIntercept(func() {
		if len(AppStatus.openNotes) > 0 {
			fmt.Println(fmt.Sprintf("len of opne notes array is %d", len(AppStatus.openNotes)))
			//do not close if there are notes open
			dlg := dialog.NewInformation("Error", "There are notes open, please close them before closing the application!", mainWindow)
			dlg.Show()
		} else {
			mainWindow.Close()
		}
	})

	mainWindow.ShowAndRun()
}

func CreateMainPanel() *fyne.Container {

	mainGridContainer := container.NewScroll(AppContainers.grid)
	AppContainers.mainGridContainer = mainGridContainer
	mainPageContainer := container.NewScroll(AppContainers.singleNoteStack)
	AppContainers.mainPageContainer = mainPageContainer
	bgRect := canvas.NewRectangle(AppTheme.MainBgColour)

	mainStackedContainer := container.NewStack(bgRect, mainPageContainer, mainGridContainer)

	return mainStackedContainer
}

func CreateTopPanel() *fyne.Container {
	//AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	spacerLabel := widget.NewLabel("                                ")
	AppWidgets.viewLabel = widget.NewLabelWithStyle("Pinned Notes      >", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	AppWidgets.pageLabel = widget.NewLabel("Page: ")
	AppWidgets.pageLabel.Hide()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.GridIcon(), func() {
			if AppStatus.currentLayout != LAYOUT_GRID {
				AppStatus.currentLayout = LAYOUT_GRID
				PageView.Reset()
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.FileIcon(), func() {
			if AppStatus.currentLayout != LAYOUT_PAGE {
				AppStatus.currentLayout = LAYOUT_PAGE
				PageView.Reset()
				UpdateView()
			}
		}),
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			if PageView.PageBack() > 0 {
				UpdateView()
			}

		}),
		widget.NewToolbarAction(theme.NavigateNextIcon(), func() {
			if PageView.PageForward() > 0 {
				UpdateView()
			}

		}),
	)

	settingsBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			fmt.Println("Setting pressed")

		}),
	)

	AppWidgets.toolbar = toolbar
	topPanel := container.New(layout.NewHBoxLayout(),
		spacerLabel,
		AppWidgets.viewLabel,
		layout.NewSpacer(),
		toolbar,
		AppWidgets.pageLabel,
		layout.NewSpacer(),
		spacerLabel,
		settingsBar,
	)

	return topPanel
}

func CreateSidePanel() *fyne.Container {
	AppContainers.searchPanel = CreateSearchPanel()
	newNoteBtn := widget.NewButtonWithIcon("+", theme.DocumentCreateIcon(), func() {
		CreateNewNote()
	})

	searchBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		//Display the search panel here
		ShowSearchPanel()
	})

	//pinnedBtn := widget.NewButton("P", func(){
	pinnedBtn := widget.NewButtonWithIcon("Pinned", theme.RadioButtonCheckedIcon(), func() {
		var err error
		AppStatus.currentView = VIEW_PINNED
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	RecentBtn := widget.NewButtonWithIcon("Recent", theme.HistoryIcon(), func() {
		var err error
		//AppStatus.notes,err = scribedb.GetRecentNotes(Conf.Settings.RecentNotesLimit)
		AppStatus.currentView = VIEW_RECENT
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	CreateNotebooksList()

	notebooksBtn := widget.NewButtonWithIcon("Notebooks", theme.FolderOpenIcon(), func() {
		ShowNotebooks()
	})

	spacerLabel := widget.NewLabel(" ")

	btnPanel := container.NewVBox(searchBtn, newNoteBtn, spacerLabel, pinnedBtn, RecentBtn, notebooksBtn)
	AppContainers.listPanel = container.NewStack(AppWidgets.notebooksList)
	AppContainers.listPanel.Hide()
	//AppContainers.searchPanel.Hide()

	sideContainer := container.NewHBox(btnPanel, AppContainers.listPanel, AppContainers.searchPanel)

	return sideContainer
}

func CreateSearchPanel() *fyne.Container {

	AppWidgets.searchResultsLabel = widget.NewLabel("")
	filterLabel := widget.NewLabel("Filter: -")
	searchFilter := widget.NewCheckGroup([]string{SEARCH_FILT_WOLE_WORDS, SEARCH_FILT_PINNED}, func(selected []string) {
		AppStatus.searchFilter.Pinned = false
		AppStatus.searchFilter.WholeWords = false
		for _, sel := range selected {
			fmt.Println("selected: " + sel)
			if sel == SEARCH_FILT_PINNED {
				AppStatus.searchFilter.Pinned = true
			}

			if sel == SEARCH_FILT_WOLE_WORDS {
				AppStatus.searchFilter.WholeWords = true
			}
		}
		var err error = UpdateView()
		if err != nil {
			log.Print("Error getting search results (after setting filter): ")
			dialog.ShowError(err, mainWindow)
			//log.Panic(err)
		}
	})
	searchLabel := widget.NewLabel("               Search:               ")
	AppWidgets.searchEntry = widget.NewEntry()
	AppWidgets.searchEntry.OnSubmitted = func(text string) {
		AppStatus.currentView = VIEW_SEARCH
		var err error = UpdateView()
		mainWindow.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		if err != nil {
			log.Print("Error getting search results: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}

	}

	searchPanel := container.NewVBox(searchLabel, AppWidgets.searchEntry, AppWidgets.searchResultsLabel, filterLabel, searchFilter)
	return searchPanel
}

func ShowNotesInGrid(notes []scribedb.NoteData, noteSize fyne.Size) {
	if AppContainers.grid == nil || AppContainers.mainGridContainer == nil {
		return
	}

	if AppContainers.mainPageContainer != nil {
		AppContainers.mainPageContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	PageView.Step = Conf.Settings.GridMaxPages
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	AppContainers.grid.RemoveAll()
	numPages := (PageView.CurrentPage + PageView.Step) - 1
	if numPages > len(notes) {
		numPages = PageView.NumberOfPages
	}

	if AppWidgets.pageLabel.Hidden != true {
		AppWidgets.pageLabel.SetText(PageView.GetGridLabelText())
	}

	for i := PageView.CurrentPage - 1; i < numPages; i++ {
		richText := NewScribeNoteText(notes[i].Content, func() {
			if slices.Contains(AppStatus.openNotes, notes[i].Id) {
				//note is already open
				fmt.Println("note is already open")
			} else {
				AppStatus.openNotes = append(AppStatus.openNotes, notes[i].Id)
				OpenNoteWindow(notes[i].Id)
			}
		})
		richText.Wrapping = fyne.TextWrapWord
		themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
		noteColour, _ := RGBStringToFyneColor(notes[i].BackgroundColour)
		noteBackground := canvas.NewRectangle(noteColour)
		if notes[i].BackgroundColour == "#e7edef" || notes[i].BackgroundColour == "#FFFFFF" {
			noteBackground = canvas.NewRectangle(AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
		}

		colourStack := container.NewStack(noteBackground)
		textPadded := container.NewPadded(themeBackground, richText)
		noteStack := container.NewStack(colourStack, textPadded)

		//borderLayout := container.NewBorder(noteBackground,noteBackground,noteBackground, noteBackground,textStack)
		AppContainers.grid.Add(noteStack)
	}
	AppContainers.grid.Refresh()
	AppContainers.mainGridContainer.Show()
}

func ShowNotesAsPages(notes []scribedb.NoteData) {
	if AppContainers.mainGridContainer != nil {
		AppContainers.mainGridContainer.Hide()
	}

	PageView.NumberOfPages = len(notes)
	PageView.Step = 1
	if PageView.CurrentPage == 0 {
		PageView.CurrentPage = 1
	}

	retreievdNote, err := scribedb.GetNote(notes[PageView.CurrentPage-1].Id)

	if err != nil {
		log.Println("error getting note")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}

	noteInfo := note.NoteInfo{
		Id:           retreievdNote.Id,
		Notebook:     retreievdNote.Notebook,
		DateCreated:  retreievdNote.Created,
		DateModified: retreievdNote.Modified,
		Colour:       retreievdNote.BackgroundColour,
		Content:      retreievdNote.Content,
		Deleted:      false,
	}

	AppWidgets.pageLabel.SetText(PageView.GetLabelText())
	AppWidgets.pageLabel.Show()

	AppWidgets.singleNotePage.ParseMarkdown(noteInfo.Content)
	AppWidgets.singleNotePage.Wrapping = fyne.TextWrapWord
	AppWidgets.singleNotePage.Refresh()

	themeBackground := canvas.NewRectangle(AppTheme.NoteBgColour)
	noteColour, _ := RGBStringToFyneColor(noteInfo.Colour)
	noteBackground := canvas.NewRectangle(noteColour)
	if noteInfo.Colour == "#e7edef" || noteInfo.Colour == "#FFFFFF" || noteInfo.Colour == "ffffff" {
		noteBackground = canvas.NewRectangle(AppTheme.NoteBgColour) // colour not set or using the old scribe default note colour
	}

	colourStack := container.NewStack(noteBackground)
	textPadded := container.NewPadded(themeBackground, AppWidgets.singleNotePage)
	noteStack := container.NewStack(colourStack, textPadded)

	AppContainers.singleNoteStack.RemoveAll()
	AppContainers.singleNoteStack.Add(noteStack)

	AppContainers.mainPageContainer.Show()
	AppContainers.mainPageContainer.Refresh()

}

func UpdateView() error {
	//var notes []scribedb.NoteData
	var err error
	//fyne.CurrentApp().SendNotification(fyne.NewNotification("Current View: ", currentView))
	switch AppStatus.currentView {
	case VIEW_PINNED:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		AppWidgets.viewLabel.SetText("Pinned Notes")
		AppStatus.notes, err = scribedb.GetPinnedNotes()
		AppStatus.currentNotebook = ""
	case VIEW_RECENT:
		if AppContainers.listPanel != nil {
			AppContainers.listPanel.Hide()
		}
		if AppContainers.searchPanel != nil {
			AppContainers.searchPanel.Hide()
		}
		AppWidgets.viewLabel.SetText(("Recent Notes"))
		AppStatus.notes, err = scribedb.GetRecentNotes(Conf.Settings.RecentNotesLimit)
		AppStatus.currentNotebook = ""
	case VIEW_NOTEBOOK:
		AppWidgets.viewLabel.SetText("Notebook - " + AppStatus.currentNotebook)
		AppStatus.notes, err = scribedb.GetNotebook(AppStatus.currentNotebook)
	case VIEW_SEARCH:
		if len(strings.TrimSpace(AppWidgets.searchEntry.Text)) > 0 {
			AppStatus.notes, err = scribedb.GetSearchResults(AppWidgets.searchEntry.Text, AppStatus.searchFilter)
			if err == nil {
				AppWidgets.searchResultsLabel.SetText(fmt.Sprintf("Found (%d) > ", len(AppStatus.notes)))
				AppWidgets.viewLabel.SetText("Search Results")
			}
		}

	default:
		err = errors.New("undefined view")
	}

	if err != nil {
		return err
	}

	switch AppStatus.currentLayout {
	case LAYOUT_GRID:
		if len(AppStatus.notes) <= Conf.Settings.GridMaxPages {
			AppWidgets.toolbar.Items[2].ToolbarObject().Hide()
			AppWidgets.toolbar.Items[3].ToolbarObject().Hide()
			AppWidgets.pageLabel.Hide()
		} else {
			AppWidgets.toolbar.Items[2].ToolbarObject().Show()
			AppWidgets.toolbar.Items[3].ToolbarObject().Show()
			AppWidgets.pageLabel.Show()
		}
		ShowNotesInGrid(AppStatus.notes, AppStatus.noteSize)
	case LAYOUT_PAGE:
		AppWidgets.toolbar.Items[2].ToolbarObject().Show()
		AppWidgets.toolbar.Items[3].ToolbarObject().Show()
		ShowNotesAsPages(AppStatus.notes)
	default:
		err = errors.New("undefined layout")
	}

	return err
}

func CreateNotebooksList() {
	var err error
	AppStatus.notebooks, err = scribedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}
	AppWidgets.notebooksList = widget.NewList(
		func() int {
			return len(AppStatus.notebooks)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("------------Notebooks (xx)------------", func() {})

		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			AppStatus.notes, _ = scribedb.GetNotebook(AppStatus.notebooks[id])
			o.(*widget.Button).SetText(fmt.Sprintf("%s (%d)", AppStatus.notebooks[id], len(AppStatus.notes)))
			o.(*widget.Button).OnTapped = func() {
				//AppStatus.notes,_ = scribedb.GetNotebook(AppStatus.notebooks[id])
				AppStatus.currentView = VIEW_NOTEBOOK
				AppStatus.currentNotebook = AppStatus.notebooks[id]
				PageView.Reset()
				UpdateView()
			}
		},
	)

}

func UpdateNotebooksList() {
	var err error
	AppStatus.notebooks, err = scribedb.GetNotebooks()
	if err != nil {
		log.Print("Error getting Notebooks: ")
		dialog.ShowError(err, mainWindow)
		log.Panic(err)
	}
	AppWidgets.notebooksList.Refresh()
}

func ShowSearchPanel() {
	if AppContainers.searchPanel.Hidden {
		AppContainers.searchPanel.Show()
		mainWindow.Canvas().Focus(AppWidgets.searchEntry)
	} else {
		mainWindow.Canvas().Unfocus() //unfocuses entry to allow keyboard shortcits ro work
		AppContainers.searchPanel.Hide()
	}
	mainWindow.Canvas().Refresh(AppContainers.searchPanel)
}

func CreateNewNote() {
	if AppContainers.listPanel != nil {
		AppContainers.listPanel.Hide()
	}
	if AppContainers.searchPanel != nil {
		AppContainers.searchPanel.Hide()
	}
	OpenNoteWindow(0) //new note has id=0
}

func ShowNotebooks() {
	UpdateNotebooksList()
	if AppStatus.currentView != VIEW_NOTEBOOK {
		AppWidgets.viewLabel.SetText("Notebooks")
	}

	if AppContainers.listPanel != nil {
		if AppContainers.listPanel.Visible() {
			AppContainers.listPanel.Hide()
		} else {
			AppContainers.listPanel.Show()
		}
	}
	PageView.Reset()
}

func AddMainKeyboardShortcuts() {
	// Add a standard shortcut (Ctrl+S)
	ctrl_p := &desktop.CustomShortcut{
		KeyName:  fyne.KeyP,
		Modifier: fyne.KeyModifierControl,
	}

	//Keyboard shortcut to show Pinned Notes
	mainWindow.Canvas().AddShortcut(ctrl_p, func(shortcut fyne.Shortcut) {
		var err error
		AppStatus.currentView = VIEW_PINNED
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting pinned notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})
	ctrl_r := &desktop.CustomShortcut{
		KeyName:  fyne.KeyR,
		Modifier: fyne.KeyModifierControl,
	}

	//Keyboard shortcut to show Recent notes
	mainWindow.Canvas().AddShortcut(ctrl_r, func(shortcut fyne.Shortcut) {
		var err error
		AppStatus.currentView = VIEW_RECENT
		PageView.Reset()
		err = UpdateView()
		if err != nil {
			log.Print("Error getting recent notes: ")
			dialog.ShowError(err, mainWindow)
			log.Panic(err)
		}
	})

	//Keyboard shortcut to show search panel
	ctrl_f := &desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: fyne.KeyModifierControl,
	}

	mainWindow.Canvas().AddShortcut(ctrl_f, func(shortcut fyne.Shortcut) {
		ShowSearchPanel()
	})

	//Keyboard shortcut to create a new note
	ctrl_shift_n := &desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
	}
	mainWindow.Canvas().AddShortcut(ctrl_shift_n, func(shortcut fyne.Shortcut) {
		CreateNewNote()
	})

	//Keyboard shortcut to show notebooks list
	ctrl_n := &desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierControl,
	}
	mainWindow.Canvas().AddShortcut(ctrl_n, func(shortcut fyne.Shortcut) {
		ShowNotebooks()
	})

}
