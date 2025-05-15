package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"scribe-nb/config"
	"scribe-nb/scribedb"
	"scribe-nb/ui"
)

const VERSION = "0.12"

func main() {
	var err error
	var dir_err error
	var appConfig *config.Config
	const confFileName = "config.toml"
	var confFilePath string
	var homeDir string

	if homeDir, dir_err = os.UserHomeDir(); dir_err != nil {
		log.Panicln(dir_err)
	}

	if runtime.GOOS == "windows" {
		//This needs to be improved but will do for now!!!!
		confFilePath = filepath.Join(homeDir, "scribe-nb")
	} else {
		confFilePath = filepath.Join(homeDir, ".config/scribe-nb") // development only
	}

	confFile := filepath.Join(confFilePath, confFileName)

	if _, f_err := os.Stat(confFile); f_err != nil {
		//write the default config confFile
		if err = os.MkdirAll(confFilePath, os.ModePerm); err != nil {
			log.Panicln("something went wrong with conf file path!!")
		}

		//create the default config.toml
		newConfig := CreateAppConfig(homeDir)

		if err = config.WriteConfig(confFile, newConfig); err != nil {
			log.Panicln(fmt.Sprintf("Error writing config file: %s", err))
		}
	}

	appConfig, err = config.GetConfig(confFile)
	if err != nil {
		log.Panicln(err)
		return
	}

	//check if the database already exists
	if _, dbf_err := os.Stat(appConfig.Settings.Database); dbf_err != nil {
		dbName := filepath.Base(appConfig.Settings.Database)
		dbPath := filepath.Dir(appConfig.Settings.Database)

		//create a new database
		if err = scribedb.CreateNew(dbName, dbPath); err != nil {
			log.Panicln(fmt.Sprintf("Something went wrong creating new db: %s", err))
		}
		scribedb.Close()
	}

	err = scribedb.Open(appConfig.Settings.Database)
	defer scribedb.Close()
	if err != nil {
		log.Panicln(err)
	}

	ui.StartUI(appConfig, confFile, VERSION)
}

func CreateAppConfig(homeDir string) config.Config {
	appSettings := config.AppSettings{
		Database: filepath.Join(homeDir, "scribe-nb", "scribeNB.db"), //this one for release
		//Database: filepath.Join(homeDir,"sync","scribe","scribeNB.db"), //temp one for dev
		InitialLayout:    "grid",
		InitialView:      "pinned",
		NoteHeight:       350,
		NoteWidth:        600,
		RecentNotesLimit: 50,
		GridMaxPages:     500,
		ThemeVariant:     "dark",
		DarkColourNote:   "#2f2f2f",
		LightColourNote:  "#e2e2e2",
		DarkColourBg:     "#1e1e1e",
		LightColourBg:    "#efefef",
	}
	newConfig := config.Config{
		Title:    fmt.Sprintf("Scribe-nb v%s", VERSION),
		Settings: appSettings,
	}

	return newConfig
}
