package scribedb

import (
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func PinNote(id uint) (int64, error) {
	var pinDate = time.Now().String()[0:19]
	res, err := db.Exec("update notes set pinned = 1, pinnedDate = ? where id = ?", pinDate, id)
	rows, _ := res.RowsAffected()
	return rows, err
}

func UnpinNote(id uint) (int64, error) {
	res, err := db.Exec("update notes set pinned = 0, pinnedDate = '' where id = ?", id)
	rows, _ := res.RowsAffected()
	return rows, err
}

func DeleteNote(id uint) (int64, error) {
	res, err := db.Exec("delete from notes where id = ?", id)
	rows, _ := res.RowsAffected()
	return rows, err
}

func SaveNote(id uint, notebook string, content string, pinned uint, pinnedDate string, colour string) (int64, error) {
	if !connected {
		return 0, errors.New("SaveNote: database not connected")
	}

	var modified = time.Now().String()[0:19]
	//fmt.Printf("fields: %d, %s, %d, %s, %s\n", id, notebook, pinned, modified, colour)

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?, pinnedDate = ?,  modified = ?, BGColour = ? where id = ?", notebook, content, pinned, pinnedDate, modified, colour, id)
	rows, _ := res.RowsAffected()

	return rows, err
}

func SaveNoteNoTimeStamp(id uint, notebook string, content string, pinned uint, pinnedDate, colour string) (int64, error) {
	if !connected {
		return 0, errors.New("SaveNote: database not connected")
	}

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?, pinnedDate = ?, BGColour = ? where id = ?", notebook, content, pinned, pinnedDate, colour, id)
	rows, _ := res.RowsAffected()

	return rows, err
}

func InsertNote(notebook string, content string, pinned uint, pinnedDate string, colour string) (int64, error) {
	if !connected {
		return 0, errors.New("InsertNote: database not connected")
	}

	//calculate date created and modified
	var created = time.Now().String()[0:19]
	var modified = created

	res, err := db.Exec("INSERT INTO notes VALUES(NULL,?,?,?,?,?,?,?)", notebook, content, created, modified, pinned, colour, pinnedDate)
	rows, _ := res.RowsAffected()

	return rows, err
}
