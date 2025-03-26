package config

import ()

type Config struct {
	Title       string
	Settings AppSettings `toml:"settings"`
}

type AppSettings struct {
	Database         string
	RecentNotesLimit int     `toml:"recent_notes_limit"`
	NoteWidth        float32 `toml:"note_width"`
	NoteHeight       float32 `toml:"note_height"`
	InitialView      string  `toml:"initial_view"`
	InitialLayout    string  `toml:"initial_layout"`
	GridMaxPages     int     `toml:"grid_max_pages"`
	DarkColour       string  `toml:"dark_colour"`
	LightColour	 string  `toml:"light_colour"`
}
