package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type scribeNoteText struct {
	widget.RichText
	OnTapped func()
}

//Implement onTapped for this widget
func (sn *scribeNoteText) Tapped(*fyne.PointEvent) {
	if sn.OnTapped != nil {
		sn.OnTapped()
	}
}

func newScribeNoteText(content string, tapped func()) *scribeNoteText {
	rt := &scribeNoteText{}
	rt.AppendMarkdown(content)
	rt.OnTapped = tapped
	return rt
}


