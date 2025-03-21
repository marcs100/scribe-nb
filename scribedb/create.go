package scribedb

import(
	"os"
	"path/filepath"
	"errors"
)

func CreateNew(dbFileName string, dbFilePath string)error{
	var err error = nil

	if err = os.MkdirAll(dbFilePath, os.ModePerm); err != nil{
		return (err)
	}

	err = Open(filepath.Join(dbFilePath, dbFileName))
	if err == nil {
		noteTableCom := "create table if not exists notes (id INTEGER PRIMARY KEY, notebook TEXT, content TEXT, created TEXT,modified TEXT, pinned INTEGER,BGColour TEXT)"
		tagTableCom := "create table if not exists tags (id INTEGER PRIMARY KEY, name TEXT, noteid INTEGER)"

		if err = createTable(noteTableCom); err != nil{
			return err
		}

		err = createTable(tagTableCom)
	}

	return err
}

func createTable(tableCom string)error{
	if db == nil{
		return errors.New("datbase not connected:  db == nil")
	}
	_, err := db.Exec(tableCom)
	return err
}
