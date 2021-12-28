package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"math"
	"strconv"
)

func CalcTotalBitrate(video *settings.Video) (float64, bool) {
	bitrate := video.Size / video.Time.Time

	bitrate = math.Min(bitrate, settings.Encoding.BitrateLimitMax)

	if bitrate < settings.Encoding.BitrateLimitMin {
		maxLength := video.Size / settings.Encoding.BitrateLimitMin
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		return 0, true
	}

	return bitrate, false
}

func CalcAudioBitrate(video *settings.Video) float64 {
	// Audio calc
    mult := 1.0
    if video.Input.AudioChannels == 1 {
    	mult = 0.5
	}
	AudioBitrate := video.Output.TotalBitrate * mult * float64(video.Output.Audio.Encoder.BitratePerc) / float64(100)

	AudioBitrate = math.Min(AudioBitrate, video.Output.Audio.Encoder.MaxBitrate * mult)
	AudioBitrate = math.Max(AudioBitrate, video.Output.Audio.Encoder.MinBitrate * mult)

	return AudioBitrate
}

func CalcOverhead(FPS float64, duration float64) float64 {
	// muxers seem to use around 300 bits per frame plus 500 bits per second, MKV/WebM is a little less but why risk it
	// calculation credit: RootAtKali
	var overhead float64
	const constantOverhead = 656.25
	const timeOverhead = 0.5
	const frameOverhead = 0.3

	overhead += constantOverhead / duration
	overhead += timeOverhead
	overhead += frameOverhead * FPS

	return overhead
}

func CalcH264Overhead(duration float64) float64 {
	return 640000 / math.Sqrt(duration)
}