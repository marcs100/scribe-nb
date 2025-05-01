package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func BackupNotes(parentWindow fyne.Window) {
	d := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err == nil {
			if dir != nil {
				//call backup funtion in scribedb here!!!!!!!!!!!!!!!
			}
		} else {
			dialog.ShowError(err, parentWindow)
		}

	}, parentWindow)
	d.Show()
}
