package vidEnc

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(video *settings.Video, pass int) (filter string, vertRes int, FPS float64) {
	var filters []string
	// FPS
	FPS = video.Input.FPS
	if float64(video.Output.Video.Limits.FPSMax) < video.Input.FPS && !settings.Original {
		if settings.Encoding.HalveDownFPS {
			for FPS > float64(video.Output.Video.Limits.FPSMax) {
				FPS /= 2
			}
		} else {
			FPS = float64(video.Output.Video.Limits.FPSMax)
		}
		filters = append(filters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	// Deduplication
	if settings.Advanced.DeduplicateFrames && !settings.Original {
		maxframes := FPS - 1
		if maxframes >= 1 {
			filters = append(filters, "mpdecimate=max=" + strconv.FormatFloat(maxframes,'f', 0, 64))
		}
	}

	// Subtitle burning
	if video.Output.Subs.SubFile != "" {
		filters = append(filters, "subtitles=si=" + strconv.Itoa(video.Input.SubtitleStream) + ":f=" + video.Output.Subs.SubFile)
	}

	vertRes = video.Input.Height
	// Resolution
	if video.Output.Video.Limits.VResMax < video.Input.Height && !settings.Original {
		vertRes = video.Output.Video.Limits.VResMax
		if pass == 1 {
			filters = append(filters, "scale=-2:" + strconv.Itoa(vertRes) + ":flags=bilinear")
		} else {
			filters = append(filters, "scale=-2:" + strconv.Itoa(vertRes) + ":flags=spline")
		}
	}

	// Pixfmt conversion
	if video.Input.Pixfmt != video.Output.Video.Encoder.Pixfmt {
		filters = append(filters, "format=" + video.Output.Video.Encoder.Pixfmt)
	}

	return strings.Join(filters, ","), vertRes, FPS
}