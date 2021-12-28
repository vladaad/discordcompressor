package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)

func SelectEncoder (video *settings.Video) *settings.Video {
	for i := range settings.Encoding.BitrateTargets {
		if settings.Encoding.BitrateTargets[i].BitrateMin < video.Size {
			target := settings.Encoding.BitrateTargets[i]
			video.Output.Video.Encoder, video.Output.Audio.Encoder = encSelect(target)
			video.Output.Video.Limits = limitSelect(target, video)
			video.Output.Video.Preset = target.Preset
			return video
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

func limitSelect(target *settings.Target, video *settings.Video) *settings.Limits {
	for i := range target.Limits {
		if target.Limits[i].Focus == settings.Focus {
			return target.Limits[i]
		}
	}
	out := target.Limits[0]
	for i := range target.Limits {
		if float64(target.Limits[i].FPSMax) >= video.Input.FPS { // If the input video is 30fps, and one of the limits is 30fps but higher res, it'll pick that one instead
			out = target.Limits[i]
		}
	}
	return out
}