package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		FFmpegExecutable:    "ffmpeg",
		FFprobeExecutable:   "ffprobe",
		QaacExecutable:      "qaac64",
	}
}

type general struct {
	FFmpegExecutable    string
	FFprobeExecutable   string
	QaacExecutable      string
}