package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
}

type OutTarget struct {
	AudioPassthrough bool
	VideoPassthrough bool
	VideoBitrate     float64
	AudioBitrate     float64
}

var OutputTarget *OutTarget
var MixTracks bool
var AudioFile string
var Debug bool
var Focus string
var Original bool
var DryRun bool