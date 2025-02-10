package scribedb

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//*********** Public functions ************************


func GetPinnedNotes() ([]NoteData, error){
	var query string = "select * from notes where pinned = 1 order by modified desc"
	return getNotes(query)
}


func GetNotebook(notebookName string) ([]NoteData, error){
	var query string = fmt.Sprintf("select * from notes where notebook = '%s'", notebookName)
	return getNotes(query)
}


func GetRecentNotes(noteCount int) ([]NoteData, error){
	var query string = fmt.Sprintf("select * from notes order by modified desc LIMIT %d", noteCount)
	return getNotes(query)
}


//************ Private functions ************************

func getNotes(query string)([]NoteData, error){
	if !connected{
		return nil, errors.New("GetPinnedNotes: database not connected")
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
