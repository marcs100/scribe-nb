package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type scribeNoteText struct {
	widget.RichText
	OnTapped func()
}

// Implement onTapped for this widget
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

type EntryCustom struct {
	widget.Entry
	onCustomShortCut func(cs *desktop.CustomShortcut)
}

func (m *EntryCustom) TypedShortcut(s fyne.Shortcut) {
	var ok bool
	var cs *desktop.CustomShortcut
	if cs, ok = s.(*desktop.CustomShortcut); !ok {
		//fmt.Printf("shortcut name is %s", cs.ShortcutName())
		m.Entry.TypedShortcut(s) //not a customshort cut - pass through to normal predifined shortcuts
		fmt.Println("** Not a custom shortcut!!")
		return
	}
	//var name = cs.ShortcutName()
	if m.onCustomShortCut == nil {
		fmt.Println("********Error nil ref to m.OnCustomShortcut()********")
		return
	}
	m.onCustomShortCut(cs)
}

func NewEntryCustom(onCustomShortcut func(cs *desktop.CustomShortcut)) *EntryCustom {
	e := &EntryCustom{}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.onCustomShortCut = onCustomShortcut
	e.ExtendBaseWidget(e)
	return e
}

type BorderCustom struct {
	//container.border
}

/*
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
}*/
