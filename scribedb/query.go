package scribedb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//*********** Public functions ************************

func GetNote(id uint) (NoteData, error) {
	var query string = fmt.Sprintf("select * from notes where id = %d", id)
	notes, err := getNotes(query)
	return notes[0], err
}

func GetPinnedNotes() ([]NoteData, error) {
	var query string = "select * from notes where pinned = 1 order by modified desc"
	return getNotes(query)
}

func GetNotebook(notebookName string) ([]NoteData, error) {
	var query string = fmt.Sprintf("select * from notes where notebook = '%s' order by modified desc", notebookName)
	return getNotes(query)
}

func GetNotebooks() ([]string, error) {
	var query string = fmt.Sprintf("select distinct notebook from notes order by notebook asc")
	return getColumn(query)
}

func CheckNotebookExists(notebook string) (bool, error) {
	var query string = fmt.Sprintf("select notebook from notes where notebook = '%s'", notebook)
	notebooks, err := getColumn(query)
	if err != nil {
		return true, err
	}

	if len(notebooks) > 0 {
		return true, err
	}

	return false, err
}

func GetRecentNotes(noteCount int) ([]NoteData, error) {
	var query string = fmt.Sprintf("select * from notes order by modified desc LIMIT %d", noteCount)
	return getNotes(query)
}

func GetSearchResults(searchQuery string, filter SearchFilter) ([]NoteData, error) {
	return getSearchResults(searchQuery, filter)
}

//************ Private functions ************************

// use to get a single column as list of strungs
func getColumn(query string) ([]string, error) {
	if !connected {
		return nil, errors.New("GetColumn: database not connected")
	}

	rows, err := db.Query(query)
	var fields []string
	defer rows.Close()
	for rows.Next() {
		var field string
		err := rows.Scan(&field)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, err
}

func getNotes(query string) ([]NoteData, error) {
	if !connected {
		return nil, errors.New("Get Notes: database not connected")
	}

	rows, err := db.Query(query)
	var notes []NoteData
	defer rows.Close()
	for rows.Next() {
		var note NoteData
		err := rows.Scan(&note.Id, &note.Notebook, &note.Content, &note.Created, &note.Modified, &note.Pinned, &note.BackgroundColour)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	//fmt.Println(notes[1].content)
	//fmt.Println(notes[2].content)
	//fmt.Println(notes[3].content)

	return notes, err
}

func getSearchResults(searchText string, filter SearchFilter) ([]NoteData, error) {
	var st1, st2, st3 string

	if !connected {
		return nil, errors.New("Get Notes: database not connected")
	}

	query := "select * from notes where content like ? "
	searchTerm := "%" + searchText + "%"

	if filter.Pinned {
		query = fmt.Sprintf("%s and pinned = 1", query)
	}

	if filter.WholeWords {
		/*
							searchQuery + ' %',
			                '% '+searchQuery+' %',
			                '% '+searchQuery+chr(10)+'%',
			                searchQuery+chr(10)+'%',
			                '%'+chr(10)+searchQuery+chr(10)+'%'
		*/
		query = "select * from notes where content like ? or content like ? or content like ? "
		st1 = fmt.Sprint("'", searchText, " %'")
		st2 = fmt.Sprint("'% ", searchText, " %'")
		st3 = fmt.Sprint("'% ", searchText, "'")
		/*searchTerm = fmt.Sprint(
		"'", searchText, " %'",
		" or content like '% ", searchText, " %'",
		" or content like '% ", searchText, "'")
		*/
		//fmt.Println(searchTerm)
		fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(st3)
	}

	query = fmt.Sprintf("%s order by modified desc", query)
	fmt.Println(query)

	var rows *sql.Rows
	var err error
	if filter.WholeWords == false {
		rows, err = db.Query(query, searchTerm)
	} else {
		rows, err = db.Query(query, st1, st2, st3)
	}
	var notes []NoteData
	defer rows.Close()
	for rows.Next() {
		var note NoteData
		err := rows.Scan(&note.Id, &note.Notebook, &note.Content, &note.Created, &note.Modified, &note.Pinned, &note.BackgroundColour)

		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	//fmt.Println(notes[1].content)
	//fmt.Println(notes[2].content)
	//fmt.Println(notes[3].content)

	return notes, err
}
