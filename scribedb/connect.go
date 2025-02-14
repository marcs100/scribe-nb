package scribedb

import(
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var connected = false
var db *sql.DB = nil

func Open() (error){
	var err error
	db, err = sql.Open("sqlite3", "/home/marc/sync/scribe/scribeNB.db")

	if err != nil {
		log.Fatal(err)
		connected = false
	}else{
		log.Println("data base open success")
		connected = true
	}

	return err
	//defer db.Close()
}

func Close(){
	if connected{
		db.Close()
	}
}
