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
		AutoUpload:        false,
		UploadService:     "catbox",
		UploadMaxMB:       100,
		UploadMaxBitrate:  3500,
		UploadEmbedLink:   true,
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
	AutoUpload        bool
	UploadService     string
	UploadMaxMB       float64
	UploadMaxBitrate  int
	UploadEmbedLink   bool
}
