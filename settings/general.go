package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		TargetSizeMB:      25,
		Hwaccel:           "auto",
		FFmpegExecutable:  "ffmpeg",
		FFprobeExecutable: "ffprobe",
		QaacExecutable:    "qaac64",
		FDKaacExecutable:  "fdkaac",
		FHGaacExecutable:  "fhgaacenc",
		OutputSuffix:      "-%smb",
	}
}

type general struct {
	TargetSizeMB      float64
	Hwaccel           string
	FFmpegExecutable  string
	FFprobeExecutable string
	QaacExecutable    string
	FDKaacExecutable  string
	FHGaacExecutable  string
	OutputSuffix      string
}
