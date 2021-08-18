package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		FFmpegExecutable:    "ffmpeg",
		FFprobeExecutable:   "ffprobe",
	}
}

type general struct {
	FFmpegExecutable    string
	FFprobeExecutable   string
}