package settings

var Advanced = initAdvanced()

func initAdvanced () *advanced{
	return &advanced{
		MixAudioTracks: false,
		NormalizeAudio: false,
		DeduplicateFrames: true,
		SubfinderLang:     "eng",
		CompatibleFormats: []*Format{{
			Container: "mp4",
			CompatibleVideoCodecs: []*VideoCodec{{
				Name: "h264",
				PixFmts: []string{"yuv420p"},
			}},
			CompatibleAudioCodecs: []*AudioCodec {{
				Name: "aac",
				SampleRates: []int{22050, 32000, 44100, 48000},
				AudioChannels: []int{1, 2},
			}}}, {
			Container: "webm",
			CompatibleVideoCodecs: []*VideoCodec{{
				Name: "h264",
				PixFmts: []string{"yuv420p"},
			}, {
				Name: "vp9",
				PixFmts: []string{"yuv420p", "yuv420p10le"},
			}},
			CompatibleAudioCodecs: []*AudioCodec {{
				Name: "aac",
				SampleRates: []int{22050, 32000, 44100, 48000},
				AudioChannels: []int{1, 2},
			}, {
				Name: "opus",
				SampleRates: []int{22050, 32000, 44100, 48000},
				AudioChannels: []int{1, 2},
			}},
		},
		},
	}
}

type advanced struct {
	MixAudioTracks    bool
	NormalizeAudio    bool
	DeduplicateFrames bool
	SubfinderLang     string
	CompatibleFormats []*Format
}

type Format struct {
	Container             string
	CompatibleVideoCodecs []*VideoCodec
	CompatibleAudioCodecs []*AudioCodec
}

type VideoCodec struct {
	Name    string
	PixFmts []string
}

type AudioCodec struct {
	Name           string
	SampleRates    []int
	AudioChannels  []int
}