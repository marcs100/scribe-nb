package ui

import (
	"scribe-nb/config"
)

var AppContainers ApplicationContainers //structure containing pointers to fyne containers
var AppWidgets ApplicationWidgets //structure containing pointers to fyne widgets
var PageView PageViewStatus // structure to track page numbers
var AppStatus ApplicationStatus // structure contauining various aoo status
var Conf *config.Config

