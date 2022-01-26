package settings

var Encoding = initEncoding()

func initEncoding() *encoding {
	return &encoding{
		HalveDownFPS:      false,
		SizeTargetMB:      8,
		BitrateTargetMult: 1,
		BitrateLimitMax:   10000,
		BitrateLimitMin:   500,
		BitrateTargets: []*Target{{
			Limits: []*Limits{
				{
					Focus:   "framerate",
					VResMax: 1080,
					FPSMax:  60,
				},
				{
					Focus:   "resolution",
					VResMax: 1080,
					FPSMax:  30,
				},
			},
			BitrateMin:   6500,
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
					VResMax: 1080,
					FPSMax:  30,
				},
			},
			BitrateMin:   3800,
			Encoder:      "x264",
			AudioEncoder: "aac",
			Preset:       "medium",
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
			Preset:       "slow",
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
			Preset:       "slow",
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
			Preset:       "slower",
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
			BitrateMin:   750,
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
			BitrateMin:   500,
			Encoder:      "x264",
			AudioEncoder: "aac",
			Preset:       "veryslow",
		},
		},
		Encoders: []*Encoder{
			{
				Name:      "x264",
				Encoder:   "libx264",
				CodecName: "h264",
				Pixfmt:    "yuv420p",
				Options:   "-aq-mode 2",
				Keyint:    10,
				PresetCmd: "-preset",
				TwoPass:   true,
				PassCmd:   "-pass",
				Container: "mp4",
			},
		},
		AudioEncoders: []*AudioEncoder{
			nil,
		},
	}
}

type encoding struct {
	HalveDownFPS      bool
	SizeTargetMB      float64
	BitrateTargetMult float64
	BitrateLimitMax   float64
	BitrateLimitMin   float64
	BitrateTargets    []*Target
	Encoders          []*Encoder
	AudioEncoders     []*AudioEncoder
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
	Name      string
	Encoder   string
	CodecName string
	Pixfmt    string
	Options   string
	Keyint    int
	PresetCmd string
	TwoPass   bool
	PassCmd   string
	Container string
}

type AudioEncoder struct {
	Name        string
	Type        string
	Encoder     string
	CodecName   string
	Options     string
	UsesBitrate bool
	MaxBitrate  float64
	MinBitrate  float64
	BitratePerc int
}
