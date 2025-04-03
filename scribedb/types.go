package scribedb

type NoteData struct {
	Id               uint
	Notebook         string
	Content          string
	Created          string
	Modified         string
	Pinned           uint
	BackgroundColour string
}

type SearchFilter struct {
	Pinned     bool
	WholeWords bool
}
