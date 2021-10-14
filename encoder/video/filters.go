package video

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(pass int, videoStats *metadata.VidStats, limit *settings.Limits, pixfmt string) string {
	var filters []string
	// FPS
	FPS = videoStats.FPS
	if float64(limit.FPSMax) < videoStats.FPS {
		if settings.Encoding.HalveDownFPS {
			for FPS > float64(limit.FPSMax) {
				FPS /= 2
			}
		} else {
			FPS = float64(limit.FPSMax)
		}
		filters = append(filters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	if settings.Advanced.DeduplicateFrames {
		maxframes := FPS - 1
		if maxframes >= 1 {
			filters = append(filters, "mpdecimate=max=" + strconv.FormatFloat(maxframes,'f', 0, 64))
		}
	}

	// Resolution
	if limit.VResMax < videoStats.Height {
		if pass == 1 {
			filters = append(filters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=bilinear")
		} else {
			filters = append(filters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=lanczos")
		}
	}

	// Yuv420p conversion
	if videoStats.Pixfmt != pixfmt {
		filters = append(filters, "format=" + pixfmt)
	}

	return strings.Join(filters, ",")
}