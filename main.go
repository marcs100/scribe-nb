package main

import (
	"log"
	"runtime"
	"scribe-nb/config"
	"scribe-nb/scribedb"
	"scribe-nb/ui"
)

func main() {
	var err error
	var appConfig *config.Config
	var confFile string
	const VERSION = "0.2"

	if runtime.GOOS == "windows"{
		//This needs to be improved but will do for now!!!!
		confFile = "C:\\users\\vboxuser\\scribe-nb\\config.toml"
	}else{
		confFile = "/home/marc/.config/scribe-nb/config_dev.toml" // development only
		//confFile = "/home/marc/.config/scribe-nb/config.toml" //release version
	}
	appConfig,err = config.GetConfig(confFile)
	if err != nil{
		log.Panicln(err)
		return
	}

	err = scribedb.Open(appConfig.AppSettings.Database)
	defer scribedb.Close()
	if err != nil{
		log.Panicln(err)
	}

	ui.StartUI(appConfig, VERSION)
}
