package video

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(pass int, videoStats *metadata.VidStats, limit *settings.Limits, pixfmt string) string {
	var filters []string
	var fpsfilters []string // scale,tmix,fps is faster than tmix,fps,scale
	var resfilters []string
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
		fpsfilters = append(fpsfilters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	// Resolution
	if limit.VResMax < videoStats.Height {
		if pass == 1 {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=bilinear")
		} else {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(limit.VResMax) + ":flags=lanczos")
		}
	}

	if !settings.Original {
		filters = append(filters, fpsfilters...)
		filters = append(filters, resfilters...)
	}

	// Yuv420p conversion
	if videoStats.Pixfmt != pixfmt {
		filters = append(filters, "format=" + pixfmt)
	}

	return strings.Join(filters, ",")
}