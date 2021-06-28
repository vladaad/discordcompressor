package encoder

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
		if settings.Encoding.TmixDownFPS && pass != 1 {
			frames := settings.VideoStats.FPS / FPS
			if frames < 2 {frames = 2}
			fpsfilters = append(fpsfilters, "tmix=frames=" + strconv.FormatFloat(frames, 'f', 0, 64))
		}
		fpsfilters = append(fpsfilters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	// Resolution
	if settings.SelectedLimits.VResMax < settings.VideoStats.Height {
		if pass == 1 {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax))
		} else {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax) + ":flags=lanczos")
		}
	}

	if settings.Encoding.TmixDownFPS && pass != 1 {
		filters = append(filters, resfilters...)
		filters = append(filters, fpsfilters...)
	} else {
		filters = append(filters, fpsfilters...)
		filters = append(filters, resfilters...)
	}

	// Yuv420p conversion
	if settings.VideoStats.Pixfmt != settings.SelectedVEncoder.Pixfmt {
		filters = append(filters, "format=" + settings.SelectedVEncoder.Pixfmt)
	}

	return strings.Join(filters, ",")
}