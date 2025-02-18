package note

import (
	"crypto/sha512"
	"fmt"
	"scribe-nb/scribedb"
)

func UpdateHash(ni *NoteInfo){
	ni.Hash = calcHash(ni.Content)
}

func calcHash(content string) string{
	hasher := sha512.New()
	hasher.Write([]byte(content))
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func CheckChanges(orig_db *scribedb.NoteData, currentNote *NoteInfo)(bool){
	if currentNote.Deleted{
		return true
	}

	if orig_db.Pinned > 0 &&  !currentNote.Pinned{
		return true
	}else if orig_db.Pinned == 0 && currentNote.Pinned{
		return true
	}

	if orig_db.BackgroundColour != currentNote.Colour{
		return true
	}

	if orig_db.Notebook != currentNote.Notebook{
		return true
	}

	//recalculate the hash for changes
	if currentNote.Hash != calcHash(currentNote.Content){
		return true
	}

	return false
}

func SaveNote(note *NoteInfo)(int64, error){
	var pinned uint = 0
	if note.Pinned {pinned = 1}
	var res int64
	var err error
	if note.Deleted{
		return 1,err //dont save note that has been deleted from the database
	}

	if note.NewNote{
		res, err = scribedb.InsertNote(note.Notebook, note.Content, pinned, note.Colour)
	}else{
		res,err = scribedb.SaveNote(note.Id, note.Notebook, note.Content, pinned, note.Colour)
	}

	return res, err
}
