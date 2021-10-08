package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"strconv"
)

func CalcTotalBitrate(size float64, duration float64) (float64, bool) {
	bitrate := size / duration
	if bitrate > settings.Encoding.BitrateLimitMax {
		bitrate = settings.Encoding.BitrateLimitMax
	}
	if bitrate < settings.Encoding.BitrateLimitMin {
		maxLength := size / settings.Encoding.BitrateLimitMin
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		return 0, true
	}
	return bitrate, false
}

func CalcAudioBitrate(targetBitrate float64, encoder settings.AudioEncoder) float64 {
	// Audio calc
	AudioBitrate := targetBitrate * float64(encoder.BitratePerc) / float64(100)
	if AudioBitrate > encoder.MaxBitrate {
		AudioBitrate = encoder.MaxBitrate
	}
	if AudioBitrate < encoder.MinBitrate {
		AudioBitrate = encoder.MinBitrate
	}
	return AudioBitrate
}