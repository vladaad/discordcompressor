package settings

import (
	"github.com/vladaad/discordcompressor/hardware"
)

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
		Scaler:      genScaler(),
		Encoders: []*Encoder{
			nil,
		},
		AEncoders: []*AudioEncoder{
			nil,
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

func genScaler() string {
	if hardware.CudaCheck("ffmpeg") == nil {
		return "cuda"
	} else {
		return "bicubic"
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
