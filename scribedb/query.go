package scribedb

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//*********** Public functions ************************

func GetNote(id uint)(NoteData, error){
	var query string = fmt.Sprintf("select * from notes where id = %d", id)
	notes,err := getNotes(query)
	return notes[0], err
}

func GetPinnedNotes() ([]NoteData, error){
	var query string = "select * from notes where pinned = 1 order by modified desc"
	return getNotes(query)
}


func GetNotebook(notebookName string) ([]NoteData, error){
	var query string = fmt.Sprintf("select * from notes where notebook = '%s'", notebookName)
	return getNotes(query)
}

func GetNotebooks() ([]string, error){
	var query string = fmt.Sprintf("select distinct notebook from notes order by id desc")
	return getColumn(query)
}


func GetRecentNotes(noteCount int) ([]NoteData, error){
	var query string = fmt.Sprintf("select * from notes order by modified desc LIMIT %d", noteCount)
	return getNotes(query)
}

func GetNotebookCovers()([]string,error){
	var query string = "select colour from notebookCovers"
	return getColumn(query)
}


//************ Private functions ************************

//use to get a single column as list of strungs
func getColumn(query string)([]string, error){
	if !connected{
		return nil, errors.New("GetColumn: database not connected")
	}

	rows, err := db.Query(query)
	var fields []string
	defer rows.Close()
	for rows.Next(){
		var field string
		err := rows.Scan(&field)

		if err != nil{
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, err
}


func getNotes(query string)([]NoteData, error){
	if !connected{
		return nil, errors.New("Get Notes: database not connected")
	}

	rows, err := db.Query(query)
	var notes []NoteData
	defer rows.Close()
	for rows.Next(){
		var note NoteData
		err := rows.Scan(&note.Id, &note.Notebook, &note.Content, &note.Created, &note.Modified, &note.Pinned, &note.BackgroundColour)

		if err != nil{
			return nil, err
		}

		notes = append(notes, note)
	}

	//fmt.Println(notes[1].content)
	//fmt.Println(notes[2].content)
	//fmt.Println(notes[3].content)

	return notes, err
}
