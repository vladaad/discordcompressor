package settings

var Debug bool
var MixAudio bool

type fileformat struct {
	General  *general
	Encoding *encoding
}

type Vid struct {
	UUID    string
	Input   *InputStats
	Output  *Out
	InFile  string
	OutFile string
	Time    *Time
}

type Out struct {
	FPS       *FPS
	Force     *Force
	Encoder   *Encoder
	AEncoder  *AudioEncoder
	Settings  *Limit
	Bitrate   *Bitrates
	AudioFile string
}

type Time struct {
	Start    float64
	Duration float64
}

type InputStats struct {
	Width     int
	Height    int
	FPS       *FPS
	Bitrate   *Bitrates
	Duration  float64
	ACodec    string
	VCodec    string
	Pixfmt    string
	ATracks   int
	AChannels int
}

type Bitrates struct {
	Total int
	Audio int
	Video int
}

type FPS struct {
	N int
	D int
}

type Force struct {
	Video     string
	Audio     string
	Container string
}
