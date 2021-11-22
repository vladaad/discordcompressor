package video

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(pass int, videoStats *metadata.VidStats, limit *settings.Limits, pixfmt string, subFilename string, streamIndex int) (filter string, vertRes int, FPS float64) {
	var filters []string
	// FPS
	FPS = videoStats.FPS
	if float64(limit.FPSMax) < videoStats.FPS && !settings.Original {
		if settings.Encoding.HalveDownFPS {
			for FPS > float64(limit.FPSMax) {
				FPS /= 2
			}
		} else {
			FPS = float64(limit.FPSMax)
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
	if subFilename != "" {
		filters = append(filters, "subtitles=si=" + strconv.Itoa(streamIndex) + ":f=" + subFilename)
	}

	vertRes = videoStats.Height
	// Resolution
	if limit.VResMax < videoStats.Height && !settings.Original {
		vertRes = limit.VResMax
		if pass == 1 {
			filters = append(filters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=bilinear")
		} else {
			filters = append(filters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=lanczos")
		}
	}

	// Pixfmt conversion
	if videoStats.Pixfmt != pixfmt {
		filters = append(filters, "format=" + pixfmt)
	}

	return strings.Join(filters, ","), vertRes, FPS
}