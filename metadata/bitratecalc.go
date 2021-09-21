package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"strconv"
)

func CalcTotalBitrate(size float64) float64 {
	bitrate := size / settings.VideoStats.Duration
	if bitrate > settings.Encoding.BitrateLimitMax {
		bitrate = settings.Encoding.BitrateLimitMax
	}
	if bitrate < settings.Encoding.BitrateLimitMin {
		maxLength := size / settings.Encoding.BitrateLimitMin
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		os.Exit(0)
	}
	return bitrate
}

func CalcAudioBitrate(targetBitrate float64) float64 {
	// Audio calc
	AudioBitrate := targetBitrate * float64(settings.SelectedAEncoder.BitratePerc) / float64(100)
	if AudioBitrate > settings.SelectedAEncoder.MaxBitrate {
		AudioBitrate = settings.SelectedAEncoder.MaxBitrate
	}
	if AudioBitrate < settings.SelectedAEncoder.MinBitrate {
		AudioBitrate = settings.SelectedAEncoder.MinBitrate
	}
	return AudioBitrate
}