package metadata

import (
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"strconv"
)

func AddOutputMetadata(video *settings.Vid) []string {
	var options []string

	options = append(options, "-metadata:s:v:0", "title="+"Video compressed using DiscordCompressor "+build.VERSION+" | "+generateVideoDescription(video))
	if video.Output.Bitrate.Audio > 0 {
		options = append(options, "-metadata:s:a:0", "title="+generateAudioDescription(video))
	}

	return options
}

func generateVideoDescription(video *settings.Vid) (description string) {
	const FPSPrecision = 0

	description = description + strconv.Itoa(video.Output.Settings.MaxVRes) + "p"
	description = description + strconv.FormatFloat(float64(video.Output.FPS.N)/float64(video.Output.FPS.D), 'f', FPSPrecision, 64) + " "
	description = description + strconv.FormatFloat(float64(video.Output.Bitrate.Video)/1024, 'f', 0, 64) + "kbit "
	description = description + " | "
	description = description + video.Output.Encoder.Args

	return description
}

func generateAudioDescription(video *settings.Vid) (description string) {
	if video.Output.APassthrough {
		description = description + "Passed through - "
		description = description + strconv.FormatFloat(float64(video.Input.Bitrate.Audio)/1024, 'f', 0, 64) + "kbit "
		description = description + video.Input.ACodec
	} else {
		description = description + strconv.FormatFloat(float64(video.Output.Bitrate.Audio)/1024, 'f', 0, 64) + "kbit "
		var encoderName string
		switch video.Output.AEncoder.Type {
		case "ffmpeg":
			encoderName = utils.GetArg(video.Output.AEncoder.Args, "-c:a")
		default:
			encoderName = video.Output.AEncoder.Type
		}
		if video.Output.AEncoder.Type == "ffmpeg" && encoderName == "aac" {
			encoderName = "FFmpeg AAC"
		}
		description = description + encoderName
	}

	return description
}
