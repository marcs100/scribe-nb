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

var AppContainers ApplicationContainers //structure containing pointers to fyne containers
var AppWidgets ApplicationWidgets       //structure containing pointers to fyne widgets
var PageView PageViewStatus             // structure to track page numbers
var AppStatus ApplicationStatus         // structure containing various app status
var Conf *config.Config
var mainApp fyne.App
var mainWindow fyne.Window
var AppTheme AppColours
