package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
}

var MixTracks bool
var AudioFile string
var Debug bool
var Focus string
var Original bool
var DryRun bool