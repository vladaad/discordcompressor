package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)

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