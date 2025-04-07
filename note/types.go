package note

type NoteInfo struct {
	Id           uint
	NewNote      bool
	Notebook     string
	Pinned       bool
	PinnedDate   string
	Colour       string
	DateCreated  string
	DateModified string
	Content      string
	Hash         string
	Deleted      bool
}

type NoteChanges struct {
	ContentChanged   bool
	ParamsChanged    bool
	PinStatusChanged bool
}
