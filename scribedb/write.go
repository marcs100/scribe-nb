package scribedb

import(
	"errors"
	//"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"


)

func PinNote(id uint)(int64, error){
	res, err := db.Exec("update notes set pinned = 1 where id = ?",id)
	rows,_ := res.RowsAffected()
	return rows, err
}

func UnpinNote(id uint)(int64, error){
	res, err := db.Exec("update notes set pinned = 0 where id = ?",id)
	rows,_ := res.RowsAffected()
	return rows, err
}

func DeleteNote(id uint)(int64, error){
	res, err := db.Exec("delete from notes where id = ?",id)
	rows,_ := res.RowsAffected()
	return rows, err
}


func SaveNote(id uint, notebook string, content string, pinned uint, colour string)(int64, error){
	if !connected{
		return 0, errors.New("SaveNote: database not connected")
	}

	var modified = time.Now().String()[0:19]
	//fmt.Printf("fields: %d, %s, %d, %s, %s\n", id, notebook, pinned, modified, colour)

	res, err := db.Exec("update notes set notebook = ?, content = ?, pinned = ?,  modified = ?, BGColour = ? where id = ?",notebook, content, pinned, modified, colour, id)
	rows,_ := res.RowsAffected()

	return rows, err
}

func InsertNote(notebook string, content string, pinned uint, colour string)(int64,error){
	if !connected{
		return 0,errors.New("InsertNote: database not connected")
	}

	//calculate date created and modified
	var created = time.Now().String()[0:19]
	var modified = created

	res, err := db.Exec("INSERT INTO notes VALUES(NULL,?,?,?,?,?,?);",notebook, content, pinned, created, modified, colour)
	rows,_ := res.RowsAffected()

	return rows, err
}
