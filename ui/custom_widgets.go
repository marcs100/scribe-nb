package ui

import (
	"log"

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
	custShortcut *desktop.CustomShortcut
	onShortCut   func()
}

func (m *EntryCustom) TypedShortcut(s fyne.Shortcut) {
	var cs *desktop.CustomShortcut
	var ok bool
	if cs, ok = s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		log.Println("error reveivng shortcut")
		return
	}

	if cs.ShortcutName() == m.custShortcut.ShortcutName() {
		if m.onShortCut != nil {
			m.onShortCut()
		} else {
			log.Println("could not set onShortcut")
		}
		return
	}
}

func NewEntryCustom(custShortcut *desktop.CustomShortcut, onShortcut func()) *EntryCustom {
	e := &EntryCustom{}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.ExtendBaseWidget(e)
	e.custShortcut = custShortcut
	e.onShortCut = onShortcut
	return e
}
