package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
}

type VidStats struct {
	Height	     int
	FPS		     float64
	Bitrate      int
	Duration     float64
	Pixfmt       string
	AudioTracks  int
	AudioCodec   string
	AudioBitrate int
	VideoCodec   string
	VideoBitrate int
}

type OutTarget struct {
	AudioPassthrough bool
	VideoPassthrough bool
	VideoBitrate     int
	AudioBitrate     int
}

var Starttime float64
var Time float64
var SelectedSettings *Target
var SelectedLimits *Limits
var SelectedVEncoder *Encoder
var SelectedAEncoder *AudioEncoder
var MaxTotalBitrate int
var VideoStats *VidStats
var OutputTarget *OutTarget
var MixTracks bool
var InputVideo string
var Debug bool
var Focus string
var Original bool