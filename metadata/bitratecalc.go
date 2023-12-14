package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"math"
)

func CalcOverhead(video *settings.Vid) *settings.Vid {
	// thanks RootAtKali, calculations taken from discordify.sh
	var overhead float64
	var header float64
	var frameOverhead float64
	var timeOverhead float64
	var marginBase float64
	mult := 1.0
	switch utils.GetArg(video.Output.Encoder.Args, "-c:v") {
	case "h264_nvenc", "h264_qsv", "h264_amf":
		header = 15000
		frameOverhead = 300
		timeOverhead = 2700
		marginBase = 720000
		mult = 0.9 // hwenc is 1-pass only, forced multiplier
	case "libx264":
		header = 15000
		frameOverhead = 300
		timeOverhead = 2700
		marginBase = 720000
		mult = 0.98 // x264 is so inaccurate in my testing this is needed for many types of footage
	case "libvpx-vp9":
		header = 9152
		frameOverhead = 60
		timeOverhead = 2680
		marginBase = 160000
	case "libaom-av1":
		header = 9152
		frameOverhead = 56
		timeOverhead = 2704
		marginBase = 320000
	default:
		log.Println("Encoder not recognized, overhead estimation may not be accurate")
		header = 12000
		frameOverhead = 178
		timeOverhead = 2704
		marginBase = 640000
	}

	extraMargin := marginBase / math.Sqrt(video.Time.Duration)

	overhead += header
	overhead += timeOverhead * video.Time.Duration
	overhead += frameOverhead * math.Min(float64(video.Input.FPS.N)/float64(video.Input.FPS.D), float64(video.Output.Settings.MaxFPS))
	overhead += extraMargin
	overhead /= video.Time.Duration

	video.Output.Bitrate.Total -= int(overhead)
	video.Output.Bitrate.Total = int(float64(video.Output.Bitrate.Total) * mult)
	return video
}

func CalcAudioBitrate(video *settings.Vid) *settings.Vid {
	// thanks RootAtKali, calculations taken from discordify.sh
	abr := int((318000 / (1 + math.Exp(-0.0000014*float64(video.Output.Bitrate.Total)))) - 154000)
	mult := video.Output.AEncoder.BMult
	if video.Input.AChannels == 1 {
		mult *= 0.5 // halve audio bitrate if mono
	}

	video.Output.Bitrate.Audio = int(float64(abr) * mult)
	// cap audio bitrate, spaghetti
	video.Output.Bitrate.Audio = int(math.Max(math.Min(float64(video.Output.Bitrate.Audio), float64(video.Output.AEncoder.BMax*1024)), float64(video.Output.AEncoder.BMin*1024)))
	video.Output.Bitrate.Video = video.Output.Bitrate.Total - abr
	return video
}
