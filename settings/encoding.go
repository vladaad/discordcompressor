package settings

var Encoding = initEncoding()

func initEncoding() *encoding {
	return &encoding{
		MaxBitrate:  15000,
		FastMode:    false,
		AutoRes:     true,
		AutoResCap:  1080,
		HalveFPS:    true,
		VariableFPS: false,
		Passthrough: true,
		ForceGetABR: true,
		Scaler:      "bicubic",
		Encoders: []*Encoder{
			{
				Name:   "fast",
				Passes: 2,
				Keyint: 10,
				Pixfmt: "yuv420p",
				Args:   "-c:v libx264 -preset medium -aq-mode 3",
			},
			{
				Name:   "normal",
				Passes: 2,
				Keyint: 10,
				Pixfmt: "yuv420p",
				Args:   "-c:v libx264 -preset slow -aq-mode 3",
			},
			{
				Name:   "slow",
				Passes: 2,
				Keyint: 10,
				Pixfmt: "yuv420p",
				Args:   "-c:v libx264 -preset veryslow -aq-mode 3",
			},
			{
				Name:   "ultra",
				Passes: 2,
				Keyint: 15,
				Pixfmt: "yuv420p10le",
				Args:   "-c:v libvpx-vp9 -lag-in-frames 25 -cpu-used 4 -auto-alt-ref 1 -arnr-maxframes 7 -arnr-strength 4 -aq-mode 0 -enable-tpl 1 -row-mt 1", // credit: BlueSwordM
			},
		},
		AEncoders: []*AudioEncoder{
			{
				Name:  "aac",
				Type:  "ffmpeg",
				BMult: 1.3,
				BMax:  192,
				BMin:  80,
				TVBR:  false,
				Args:  "-c:a aac -aac_coder twoloop",
			},
			{
				Name:  "opus",
				Type:  "ffmpeg",
				BMult: 1,
				BMax:  128,
				BMin:  32,
				TVBR:  false,
				Args:  "-c:a libopus",
			},
		},
		Limits: []*Limit{
			{
				Encoder:    "fast",
				AEncoder:   "aac",
				Container:  "mp4",
				MinBitrate: 12000,
				MaxVRes:    1080,
				MaxFPS:     60,
			},
			{
				Encoder:    "normal",
				AEncoder:   "aac",
				Container:  "mp4",
				MinBitrate: 5000,
				MaxVRes:    1080,
				MaxFPS:     60,
			},
			{
				Encoder:    "normal",
				AEncoder:   "aac",
				Container:  "mp4",
				MinBitrate: 3000,
				MaxVRes:    900,
				MaxFPS:     60,
			},
			{
				Encoder:    "normal",
				AEncoder:   "aac",
				Container:  "mp4",
				MinBitrate: 1500,
				MaxVRes:    720,
				MaxFPS:     60,
			},
			{
				Encoder:    "slow",
				AEncoder:   "aac",
				Container:  "mp4",
				MinBitrate: 800,
				MaxVRes:    540,
				MaxFPS:     60,
			},
			{
				Encoder:    "ultra",
				AEncoder:   "opus",
				Container:  "webm",
				MinBitrate: 0,
				MaxVRes:    540,
				MaxFPS:     30,
			},
		},
	}
}

type encoding struct {
	MaxBitrate  int
	FastMode    bool
	AutoRes     bool
	AutoResCap  int
	HalveFPS    bool
	VariableFPS bool
	Scaler      string
	Passthrough bool
	ForceGetABR bool
	Encoders    []*Encoder
	AEncoders   []*AudioEncoder
	Limits      []*Limit
}

type Encoder struct {
	Name   string
	Passes int
	Keyint int
	Pixfmt string
	Args   string
}

type AudioEncoder struct {
	Name  string
	Type  string
	BMult float64
	BMax  int
	BMin  int
	TVBR  bool
	Args  string
}

type Limit struct {
	Encoder    string
	AEncoder   string
	Container  string
	MinBitrate int
	MaxVRes    int
	MaxFPS     int
}
