package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		FFmpegExecutable:    "ffmpeg",
		FFprobeExecutable:   "ffprobe",
		QaacExecutable:      "qaac64",
		BatchModeThreads:    1,
	}
}

type general struct {
	FFmpegExecutable    string
	FFprobeExecutable   string
	QaacExecutable      string
	BatchModeThreads    int
}