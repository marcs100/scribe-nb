package scribedb

import(
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"


)


func SaveNote(id uint, notebook string, content string, pinned uint, colour string)(int64, error){
	if !connected{
		return 0, errors.New("SaveNote: database not connected")
	}

	var modified = "2025-02-17 18:24:09"

	fmt.Printf("fields: %d, %s, %d, %s, %s\n", id, notebook, pinned, modified, colour)

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?,  modified = ?, BGColour = ? where id = ?",notebook, content, pinned, modified, colour, id)
	rows,_ := res.RowsAffected()

	return rows, err
}

func InsertNote(notebook string, content string, pinned uint, colour string)(int64,error){
	if !connected{
		return 0,errors.New("InsertNote: database not connected")
	}

	//calculate date created and modified
	var created string = "2025-02-17 17:53:03"
	var modified = created

	res, err := db.Exec("INSERT INTO notes VALUES(NULL,?,?,?,?,?,?);",notebook, content, pinned, created, modified, colour)
	rows,_ := res.RowsAffected()

	return rows, err
}
