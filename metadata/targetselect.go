package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)

func SelectEncoder (bitrate float64) (*settings.Encoder, *settings.AudioEncoder, *settings.Target, *settings.Limits) {
	for i := range settings.Encoding.BitrateTargets {
		if settings.Encoding.BitrateTargets[i].BitrateMin < bitrate {
			target := settings.Encoding.BitrateTargets[i]
			venc, aenc := encSelect(target)
			limits := limitSelect(target)
			return venc, aenc, target, limits
		}
	}
	panic("Could not find suitable bitrate target")
}

func encSelect (target *settings.Target) (*settings.Encoder, *settings.AudioEncoder) {
	for i := range settings.Encoding.Encoders {
		if settings.Encoding.Encoders[i].Name == target.Encoder {
			venc := settings.Encoding.Encoders[i]
			aenc := aencSelect(target.AudioEncoder)
			return venc, aenc
		}
	}
	panic("Could not find video encoder " + target.Encoder)
}

func aencSelect (encoderName string) *settings.AudioEncoder {
	for i := range settings.Encoding.AudioEncoders {
		if settings.Encoding.AudioEncoders[i].Name == encoderName {
			return settings.Encoding.AudioEncoders[i]
		}
	}
	panic("Could not find audio encoder " + encoderName)
}

func limitSelect(target *settings.Target) *settings.Limits {
	for i := range target.Limits {
		if target.Limits[i].Focus == settings.Focus {
			return target.Limits[i]
		}
	}
	return target.Limits[0]
}