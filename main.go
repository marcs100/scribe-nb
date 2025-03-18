package main

import (
	"fmt"
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
	fmt.Println("Scribe Nota Deme v0.02")

	if runtime.GOOS == "windows"{
		//This needs to be improved but will do for now!!!!
		confFile = "C:\\users\\marks\\scribe-nb\\config.toml"
	}else{
		confFile = "/home/marc/.config/scribe-nb/config_dev.toml" // development only
		//conf_file := "/home/marc/.config/scribe-nb/config.toml" //release version
	}
	appConfig,err = config.GetConfig(confFile)
	if err != nil{
		log.Panicln(err)
		return
	}

	err = scribedb.Open()
	defer scribedb.Close()
	if err != nil{
		log.Panicln(err)
	}

	ui.StartUI(appConfig)
}
