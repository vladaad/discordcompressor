package video

import (
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func addMetadata(oTarget *OutTarget, videoStats *metadata.VidStats, eOptions *settings.Encoder, vertRes int, FPS float64, eTarget *settings.Target, aOptions *settings.AudioEncoder) []string {
	var options []string

	options = append(options, "-metadata:s:v:0", "title=" + strconv.FormatFloat(settings.TargetSize, 'f', 0, 64 ) + "mb video compressed using DiscordCompressor " + build.VERSION + " | " + generateVideoDescription(oTarget, videoStats, eOptions, vertRes, FPS, eTarget))
	if oTarget.AudioBitrate != 0 {
		options = append(options, "-metadata:s:a:0", "title=" + generateAudioDescription(oTarget, videoStats, aOptions))
	}

	return options
}

func generateVideoDescription(oTarget *OutTarget, videoStats *metadata.VidStats, eOptions *settings.Encoder, vertRes int, FPS float64, eTarget *settings.Target) (description string) {
	const FPSPrecision = 0

	// Video
	if oTarget.VideoPassthrough {
		description = description + "Passed through - "
		description = description + strconv.Itoa(videoStats.Height) + "p"
		description = description + strconv.FormatFloat(videoStats.FPS, 'f', FPSPrecision, 64) + " "
		description = description + strconv.FormatFloat(videoStats.VideoBitrate, 'f', 0, 64) + "kbit "
		description = description + videoStats.VideoCodec
	} else {
		description = description + strconv.Itoa(vertRes) + "p"
		description = description + strconv.FormatFloat(FPS, 'f', FPSPrecision, 64) + " "
		description = description + strconv.FormatFloat(oTarget.VideoBitrate, 'f', 0, 64) + "kbit "
		if eOptions.TwoPass {
			description = description + "2-pass "
		} else {
			description = description + "1-pass "
		}
		description = description + eOptions.Encoder + " | "
		description = description + "-preset " + eTarget.Preset + " "
		description = description + eOptions.Options
	}

	return description
}

func generateAudioDescription(oTarget *OutTarget, videoStats *metadata.VidStats, aOptions *settings.AudioEncoder) (description string) {
	// Audio
	if oTarget.AudioPassthrough {
		description = description + "Passed through - "
		description = description + strconv.FormatFloat(videoStats.AudioBitrate, 'f', 0, 64) + "kbit "
		description = description + videoStats.AudioCodec
	} else {
		description = description + strconv.FormatFloat(oTarget.AudioBitrate, 'f', 0, 64) + "kbit "
		var encoderName string
		switch aOptions.Type {
		case "ffmpeg":
			encoderName = aOptions.Encoder
		default:
			encoderName = aOptions.Type
		}
		if aOptions.Type == "ffmpeg" && encoderName == "aac" {
			encoderName = "FFmpeg AAC"
		}
		description = description + encoderName
		if settings.Advanced.NormalizeAudio {
			description = description + " (normalized)"
		}
		description = description + " | "
		description = description + aOptions.Options
	}


	return description
}