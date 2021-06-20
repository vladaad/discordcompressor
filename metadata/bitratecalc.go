package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)


func SelectEncoder (bitrate int) bool {
	for i := range settings.Encoding.BitrateTargets {
		if settings.Encoding.BitrateTargets[i].BitrateMin < bitrate {
			settings.SelectedSettings = settings.Encoding.BitrateTargets[i]
			vencSelect(settings.SelectedSettings.Encoder)
			return true
		}
	}
	panic("Could not find suitable bitrate target")
}

func vencSelect (encoderName string) bool {
	for i := range settings.Encoding.Encoders {
		if settings.Encoding.Encoders[i].Name == encoderName {
			settings.SelectedVEncoder = settings.Encoding.Encoders[i]
			aencSelect(settings.SelectedSettings.AudioEncoder)
			return true
		}
	}
	panic("Could not find video encoder " + encoderName)
}

func aencSelect (encoderName string) bool {
	for i := range settings.Encoding.AudioEncoders {
		if settings.Encoding.AudioEncoders[i].Name == encoderName {
			settings.SelectedAEncoder = settings.Encoding.AudioEncoders[i]
			limitSelect()
			return true
		}
	}
	panic("Could not find audio encoder " + encoderName)
}

func limitSelect() bool {
	for i := range settings.SelectedSettings.Limits {
		if settings.SelectedSettings.Limits[i].Focus == settings.Focus {
			settings.SelectedLimits = settings.SelectedSettings.Limits[i]
			return true
		}
	}
	settings.SelectedLimits = settings.SelectedSettings.Limits[0]
	return false
}