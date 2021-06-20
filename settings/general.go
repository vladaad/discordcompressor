package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		FFmpegExecutable:    "ffmpeg",
		FFprobeExecutable:   "ffprobe",
		UseCustomOutputDir:  false,
		CustomOutputDir:     "",
	}
}

type general struct {
	FFmpegExecutable    string
	FFprobeExecutable   string
	UseCustomOutputDir  bool
	CustomOutputDir     string
}