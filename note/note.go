package note

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"scribe-nb/scribedb"
)

func UpdateHash(ni *NoteInfo) {
	ni.Hash = calcHash(ni.Content)
}

func calcHash(content string) string {
	hasher := sha512.New()
	hasher.Write([]byte(content))
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func CheckChanges(orig_db *scribedb.NoteData, currentNote *NoteInfo) NoteChanges {

	noteChanges := NoteChanges{}

	//recalculate the hash for changes
	if currentNote.Hash != calcHash(currentNote.Content) {
		//fmt.Println("hash mismatch!!!!!!!!") // ******* Debug only **************
		noteChanges.ContentChanged = true
	}

	if currentNote.Deleted {
		noteChanges.ContentChanged = true
		//fmt.Println("note deleted!!!!!!!!") // ******* Debug only **************
	}

	if orig_db.Pinned > 0 && !currentNote.Pinned {
		noteChanges.PinStatusChanged = true
		//fmt.Println("pinned status changed!!!!!!!!") // ******* Debug only **************
	} else if orig_db.Pinned == 0 && currentNote.Pinned {
		noteChanges.PinStatusChanged = true
	}

	if orig_db.BackgroundColour != currentNote.Colour {
		noteChanges.ParamsChanged = true
		//fmt.Println("note colour changed!!!!!!!!") // ******* Debug only **************
	}

	if orig_db.Notebook != currentNote.Notebook {
		noteChanges.ParamsChanged = true
		//fmt.Println("notebook changed!!!!!!!!") // ******* Debug only **************
	}

	return noteChanges
}

func SaveNote(note *NoteInfo) (int64, error) {
	var pinned uint = 0
	if note.Pinned {
		pinned = 1
	}
	var res int64
	var err error
	var id int64

	if note.Deleted {
		return 1, err //dont save note that has been deleted from the database
	}

	if note.NewNote {
		res, id, err = scribedb.InsertNote(note.Notebook, note.Content, pinned, note.PinnedDate, note.Colour)
		if err == nil {
			note.NewNote = false
			note.Id = uint(id)
		}
	} else {
		res, err = scribedb.SaveNote(note.Id, note.Notebook, note.Content, pinned, note.PinnedDate, note.Colour)
	}

	if err == nil {
		note.Hash = calcHash(note.Content) // update hash for saved note
	}

	return res, err
}

func SaveNoteNoTimeStamp(note *NoteInfo) (int64, error) {
	var pinned uint = 0
	if note.Pinned {
		pinned = 1
	}
	var res int64
	var err error
	if note.Deleted {
		return 1, err //dont save note that has been deleted from the database
	}

	if note.NewNote {
		res = 0
		err = errors.New("SaveNoteNoTimeStamp: error - can't save a new note in this function")
	} else {
		res, err = scribedb.SaveNoteNoTimeStamp(note.Id, note.Notebook, note.Content, pinned, note.PinnedDate, note.Colour)
	}

	return res, err
}

func GetPropertiesText(noteInfo *NoteInfo) string {
	pinnedDate := noteInfo.PinnedDate
	var pinnedStat string = "no"
	if noteInfo.Pinned {
		pinnedStat = "yes"
	} else {
		pinnedDate = "n/a"
	}
	created := noteInfo.DateCreated
	modified := noteInfo.DateModified

	if len(created) > 16 {
		created = created[:16]
	}
	if len(modified) > 16 {
		modified = modified[:16]
	}
	if len(pinnedDate) > 16 {
		pinnedDate = pinnedDate[:16]
	}
	return fmt.Sprintf("note id: %d\r\nnotebook: %s\r\n\r\ncreated:   %s\r\nmodified: %s\r\n\r\npinned: %s\r\ndate pinned: %s\r\n\r\nborder colour: %s",
		noteInfo.Id, noteInfo.Notebook, created, modified, pinnedStat, pinnedDate, noteInfo.Colour)
}
