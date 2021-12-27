package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
	Advanced  *advanced
}

type VidStats struct {
	Height	     int
	FPS		     float64
	Bitrate      float64
	Duration     float64
	Pixfmt       string
	AudioTracks  int
	AudioCodec   string
	AudioBitrate float64
	SampleRate   int
	AudioChannels int
	VideoCodec   string
	VideoBitrate float64
	MatchingSubs bool
	SubtitleStream int
}

type Video struct {
	UUID string
	Filename string
	Size float64
	Input *VidStats
	Output *Out
	Time *Time
}

type Out struct {
	Video *VideoOut
	Audio *AudioOut
	TotalBitrate float64
	Subs *SubOut
}

type VideoOut struct {
	Bitrate float64
	Encoder *Encoder
	Preset string
	Limits *Limits
	Passthrough bool
}

type AudioOut struct {
	Bitrate float64
	Encoder *AudioEncoder
	Mix bool
	Normalize bool
	Filename string
	Passthrough bool
}

type SubOut struct {
	BurnSubs bool
	SubFile string
}

type Time struct {
	Start float64
	Time float64
}

var ForceScore float64
var Debug bool
var Focus string
var Original bool
var DryRun bool
var ShowStdOut bool
var BatchMode bool
var TargetSize float64