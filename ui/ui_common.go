package ui

import (
	"scribe-nb/config"

	"fyne.io/fyne/v2"
)

const (
	VIEW_PINNED   string = "pinned"
	VIEW_RECENT   string = "recent"
	VIEW_NOTEBOOK string = "notebooks"
	VIEW_SEARCH   string = "search"
	VIEW_TAGS     string = "tag"
)

const (
	LAYOUT_GRID string = "grid"
	LAYOUT_PAGE string = "page"
)

const (
	SEARCH_FILT_PINNED     string = "Pinned"
	SEARCH_FILT_WOLE_WORDS string = "Whole words only"
)

const (
	EDIT_MODE string = "Edit"
	VIEW_MODE string = "View"
)

var mainApp fyne.App
var mainWindow fyne.Window
var noteWindow fyne.Window
var AppContainers ApplicationContainers //structure containing pointers to fyne containers for main window
var NoteContainers NoteWindowContainers //structure containing poineters to fyne containers for note window
var AppWidgets ApplicationWidgets       //structure containing pointers to fyne widgets for main window
var NoteWidgets NoteWindowWidgets       //structure containing pointers to fyne widgets for note window
var NoteCanvas NoteWindowCanvas         //structure containing pinters to canvas object for note window
var PageView PageViewStatus             //structure to track page numbers
var AppStatus ApplicationStatus         //structure containing various app status
var Conf *config.Config
var AppTheme AppColours
