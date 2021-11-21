package settings

// Stolen from https://github.com/Wieku/danser-go/app/settings

type fileformat struct {
	General   *general
	Decoding  *decoding
	Encoding  *encoding
	Advanced  *advanced
}

var ForceScore float64
var Debug bool
var Focus string
var Original bool
var DryRun bool
var ShowStdOut bool
var BatchMode bool
var TargetSize float64