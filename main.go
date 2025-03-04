package main

import (
	"fmt"
	"scribe-nb/scribedb"
	"scribe-nb/ui"
	"scribe-nb/config"
)

func main() {
	fmt.Println("Scribe Nota Deme v0.01")

	conf_file := "/home/marc/.config/scribe/config.toml"

	config.GetConfig(conf_file)


	err := scribedb.Open()
	defer scribedb.Close()
	if err != nil{
		fmt.Println("Bollocks got error")
	}

	ui.StartUI()
}
