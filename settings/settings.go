package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
}

type VidStats struct {
	Height	 int
	FPS		 float64
	Bitrate  int
	Duration float64
	Pixfmt   string
}

var Starttime float64
var Time float64
var SelectedSettings *Target
var SelectedLimits *Limits
var SelectedVEncoder *Encoder
var SelectedAEncoder *AudioEncoder
var VideoStats *VidStats
var InputVideo string
var Debug bool
var Focus string
var Original bool