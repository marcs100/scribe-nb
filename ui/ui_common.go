package ui

import (
	"scribe-nb/config"
	"fyne.io/fyne/v2"
)

const VIEW_PINNED string = "pinned"
const VIEW_RECENT string = "recent"
const VIEW_NOTEBOOK string = "notebooks"
const VIEW_SEARCH string = "search"
const VIEW_TAGS string = "tag"
const LAYOUT_GRID = "grid"
const LAYOUT_PAGE = "page"


var AppContainers ApplicationContainers //structure containing pointers to fyne containers
var AppWidgets ApplicationWidgets //structure containing pointers to fyne widgets
var PageView PageViewStatus // structure to track page numbers
var AppStatus ApplicationStatus // structure containing various app status
var Conf *config.Config
var mainApp fyne.App
var mainWindow fyne.Window


