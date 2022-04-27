package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)

func SelectEncoder(video *settings.Vid) *settings.Vid {
	for i := range settings.Encoding.Limits {
		if settings.Encoding.Limits[i].MinBitrate*1024 < video.Output.Bitrate.Total {
			target := settings.Encoding.Limits[i]
			video.Output.Settings = target
			if video.Output.Force.Container != "" {
				video.Output.Settings.Container = video.Output.Force.Container
			}
			video.Output.Encoder, video.Output.AEncoder = encSelect(video)
			return video
		}
	}
	panic("Could not find suitable bitrate target")
}

func encSelect(target *settings.Vid) (*settings.Encoder, *settings.AudioEncoder) {
	for i := range settings.Encoding.Encoders {
		if target.Output.Force.Video != "" {
			target.Output.Settings.Encoder = target.Output.Force.Video
		}
		if settings.Encoding.Encoders[i].Name == target.Output.Settings.Encoder {
			venc := settings.Encoding.Encoders[i]
			aenc := aencSelect(target)
			return venc, aenc
		}
	}
	panic("Could not find video encoder " + target.Output.Settings.Encoder)
}

func aencSelect(target *settings.Vid) *settings.AudioEncoder {
	for i := range settings.Encoding.AEncoders {
		if target.Output.Force.Audio != "" {
			target.Output.Settings.AEncoder = target.Output.Force.Audio
		}
		if settings.Encoding.AEncoders[i].Name == target.Output.Settings.AEncoder {
			return settings.Encoding.AEncoders[i]
		}
	}
	panic("Could not find audio encoder " + target.Output.Settings.AEncoder)
}
