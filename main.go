package main

import (
	"fmt"
	"scribe-nb/scribedb"
	"scribe-nb/ui"
)

func main() {
	fmt.Println("Scribe Nota Deme v0.01")
	err := scribedb.Open()
	defer scribedb.Close()
	if err != nil{
		fmt.Println("Bollocks got error")
	}

	ui.StartUI()
}
