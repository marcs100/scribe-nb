package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// View pinned notes
var ctrl_p = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl,
}

// View recent Notes
var ctrl_r = &desktop.CustomShortcut{
	KeyName:  fyne.KeyR,
	Modifier: fyne.KeyModifierControl,
}

// View Notebooks (toggles notebook list panel)
var ctrl_n = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl,
}

// open search panel
var ctrl_f = &desktop.CustomShortcut{
	KeyName:  fyne.KeyF,
	Modifier: fyne.KeyModifierControl,
}

// Open a new note
var ctrl_shift_n = &desktop.CustomShortcut{
	KeyName:  fyne.KeyN,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Set edit mode
var ctrl_shift_e = &desktop.CustomShortcut{
	KeyName:  fyne.KeyE,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Set view mode
var ctrl_shift_q = &desktop.CustomShortcut{
	KeyName:  fyne.KeyQ,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Pin note
var ctrl_shift_p = &desktop.CustomShortcut{
	KeyName:  fyne.KeyP,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Change note colour
var ctrl_shift_c = &desktop.CustomShortcut{
	KeyName:  fyne.KeyC,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}

// Open notebooks menu
var ctrl_shift_b = &desktop.CustomShortcut{
	KeyName:  fyne.KeyB,
	Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
}
