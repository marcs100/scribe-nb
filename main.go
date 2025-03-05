package main

import (
	"fmt"
	"log"
	"scribe-nb/config"
	"scribe-nb/scribedb"
	"scribe-nb/ui"
)

func main() {
	var err error
	var appConfig *config.Config
	fmt.Println("Scribe Nota Deme v0.01")

	conf_file := "/home/marc/.config/scribe/config.toml"
	appConfig,err = config.GetConfig(conf_file)
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
