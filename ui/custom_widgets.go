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

func NewScribeNoteText(content string, tapped func()) *scribeNoteText {
	rt := &scribeNoteText{}
	rt.AppendMarkdown(content)
	rt.OnTapped = tapped
	return rt
}


type buttonWithPos struct {
	widget.Button
	OnTapped func(*fyne.PointEvent)
}

//Implement onTapped with mouse position for this widget
func (bn *buttonWithPos) Tapped(e *fyne.PointEvent) {
	if bn.OnTapped != nil {
		bn.OnTapped(e)
	}
}

func NewButtonWithPos(label string, tapped func(*fyne.PointEvent)) *buttonWithPos {
	bn := &buttonWithPos{}
	bn.Text = label
	bn.OnTapped = tapped
	return bn
}


