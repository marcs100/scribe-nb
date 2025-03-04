package config

import()

type Config struct {
	Title string
	AppSettings Settings `toml:"settings"`
}


type Settings struct{
	database string
	recentNotesLimit int
	noteWidth float32
	noteHeight float32
}
