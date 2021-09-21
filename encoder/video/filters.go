package video

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(pass int) string {
	var filters []string
	var fpsfilters []string // scale,tmix,fps is faster than tmix,fps,scale
	var resfilters []string
	// FPS
	FPS = settings.VideoStats.FPS
	if float64(settings.SelectedLimits.FPSMax) < settings.VideoStats.FPS {
		if settings.Encoding.HalveDownFPS {
			for FPS > float64(settings.SelectedLimits.FPSMax) {
				FPS /= 2
			}
		} else {
			FPS = float64(settings.SelectedLimits.FPSMax)
		}
		fpsfilters = append(fpsfilters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	// Resolution
	if settings.SelectedLimits.VResMax < settings.VideoStats.Height {
		if pass == 1 {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax) + ":flags=bilinear")
		} else {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax) + ":flags=lanczos")
		}
	}

	if !settings.Original {
		filters = append(filters, fpsfilters...)
		filters = append(filters, resfilters...)
	}

	// Yuv420p conversion
	if settings.VideoStats.Pixfmt != settings.SelectedVEncoder.Pixfmt {
		filters = append(filters, "format=" + settings.SelectedVEncoder.Pixfmt)
	}

	return strings.Join(filters, ",")
}