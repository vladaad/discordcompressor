package settings

var Encoding = initEncoding()

func initEncoding() *encoding {
	return &encoding{
		TmixDownFPS:         false,
		HalveDownFPS:        false,
		SizeTargetMB:        8,
		BitrateTargetMult:   1,
		BitrateLimitMax:     12500,
		BitrateLimitMin:     500,
		BitrateTargets: []*Target{{
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 1080,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 2160,
						FPSMax:  30,
					},
				},
				BitrateMin:   7500,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "fast",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 1080,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 1440,
						FPSMax:  30,
					},
				},
				BitrateMin:   5000,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "medium",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 1080,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 1440,
						FPSMax:  30,
					},
				},
				BitrateMin:   3750,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "slow",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 900,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 1080,
						FPSMax:  30,
					},
				},
				BitrateMin:   2500,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "slower",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 720,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 1080,
						FPSMax:  30,
					},
				},
				BitrateMin:   1500,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "slower",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 720,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 900,
						FPSMax:  30,
					},
				},
				BitrateMin:   1000,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "veryslow",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 540,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 720,
						FPSMax:  30,
					},
				},
				BitrateMin:   650,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "veryslow",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 360,
						FPSMax:  60,
					},
					{
						Focus:   "resolution",
						VResMax: 540,
						FPSMax:  30,
					},
				},
				BitrateMin:   400,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "veryslow",
			}, {
				Limits: []*Limits{
					{
						Focus:   "framerate",
						VResMax: 360,
						FPSMax:  30,
					},
					{
						Focus:   "resolution",
						VResMax: 540,
						FPSMax:  15,
					},
				},
				BitrateMin:   0,
				Encoder:      "x264",
				AudioEncoder: "aac",
				Preset:       "veryslow",
			},
		},
		Encoders: []*Encoder{
			{
				Name:         "x264",
				Encoder:      "libx264",
				CodecName:    "h264",
				Pixfmt:       "yuv420p",
				Options:      "-x264-params qpmin=25",
				Keyint:       10,
				PresetCmd:    "-preset",
				TwoPass:      true,
				PassCmd:      "-pass",
				Container:    "mp4",
			},
		},
		AudioEncoders: []*AudioEncoder{
			{
				Name:         "aac",
				Type:         "ffmpeg",
				Encoder:      "aac",
				CodecName:    "aac",
				Options:      "",
				UsesBitrate:  true,
				MaxBitrate:   256,
				MinBitrate:   128,
				BitratePerc:  12,
			},
		},
	}
	}

type encoding struct {
	TmixDownFPS           bool
	HalveDownFPS          bool
	SizeTargetMB          float64
	BitrateTargetMult     float64
	BitrateLimitMax       float64
	BitrateLimitMin       float64
	BitrateTargets        []*Target
	Encoders              []*Encoder
	AudioEncoders         []*AudioEncoder
}

type Target struct {
	Limits       []*Limits
	BitrateMin   float64
	Encoder      string
	AudioEncoder string
	Preset       string
}

type Limits struct {
	Focus   string
	VResMax int
	FPSMax  int
}

type Encoder struct {
	Name         string
	Encoder      string
	CodecName    string
	Pixfmt       string
	Options      string
	Keyint       int
	PresetCmd    string
	TwoPass      bool
	PassCmd      string
	Container    string
}

type AudioEncoder struct {
	Name         string
	Type         string
	Encoder      string
	CodecName    string
	Options      string
	UsesBitrate  bool
	MaxBitrate   float64
	MinBitrate   float64
	BitratePerc  int
}