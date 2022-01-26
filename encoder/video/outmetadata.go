package vidEnc

import (
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func addMetadata(video *settings.Video, vertRes int, FPS float64) []string {
	var options []string

	options = append(options, "-metadata:s:v:0", "title="+strconv.FormatFloat(settings.TargetSize, 'f', 0, 64)+"mb video compressed using DiscordCompressor "+build.VERSION+" | "+generateVideoDescription(video, vertRes, FPS))
	if video.Output.Audio.Bitrate > 0 {
		options = append(options, "-metadata:s:a:0", "title="+generateAudioDescription(video))
	}

	return options
}

func generateVideoDescription(video *settings.Video, vertRes int, FPS float64) (description string) {
	const FPSPrecision = 0

	// Video
	if video.Output.Video.Passthrough {
		description = description + "Passed through - "
		description = description + strconv.Itoa(video.Input.Height) + "p"
		description = description + strconv.FormatFloat(video.Input.FPS, 'f', FPSPrecision, 64) + " "
		description = description + strconv.FormatFloat(video.Input.VideoBitrate, 'f', 0, 64) + "kbit "
		description = description + video.Input.VideoCodec
	} else {
		description = description + strconv.Itoa(vertRes) + "p"
		description = description + strconv.FormatFloat(FPS, 'f', FPSPrecision, 64) + " "
		description = description + strconv.FormatFloat(video.Output.Video.Bitrate, 'f', 0, 64) + "kbit "
		if video.Output.Video.Encoder.TwoPass {
			description = description + "2-pass "
		} else {
			description = description + "1-pass "
		}
		description = description + video.Output.Video.Encoder.Encoder + " | "
		description = description + "-preset " + video.Output.Video.Preset + " "
		description = description + video.Output.Video.Encoder.Options
	}

	return description
}

func generateAudioDescription(video *settings.Video) (description string) {
	// Audio
	if video.Output.Audio.Passthrough {
		description = description + "Passed through - "
		description = description + strconv.FormatFloat(video.Input.AudioBitrate, 'f', 0, 64) + "kbit "
		description = description + video.Input.AudioCodec
	} else {
		description = description + strconv.FormatFloat(video.Output.Audio.Bitrate, 'f', 0, 64) + "kbit "
		var encoderName string
		switch video.Output.Audio.Encoder.Type {
		case "ffmpeg":
			encoderName = video.Output.Audio.Encoder.Encoder
		default:
			encoderName = video.Output.Audio.Encoder.Type
		}
		if video.Output.Audio.Encoder.Type == "ffmpeg" && encoderName == "aac" {
			encoderName = "FFmpeg AAC"
		}
		description = description + encoderName
		if video.Output.Audio.Normalize {
			description = description + " (normalized)"
		}
		description = description + " | "
		description = description + video.Output.Audio.Encoder.Options
	}

	return description
}
