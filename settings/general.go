package settings

var General = initGeneral()

func initGeneral() *general {
	return &general{
		FFmpegExecutable:  "ffmpeg",
		FFprobeExecutable: "ffprobe",
		QaacExecutable:    "qaac64",
		FDKaacExecutable:  "fdkaac",
		BatchModeThreads:  1,
		OutputSuffix:      "-%smb",
	}
}

type general struct {
	FFmpegExecutable  string
	FFprobeExecutable string
	QaacExecutable    string
	FDKaacExecutable  string
	BatchModeThreads  int
	OutputSuffix      string
}
