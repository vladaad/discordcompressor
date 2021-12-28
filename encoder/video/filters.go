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

	// Resolution
	vertRes = video.Input.Height
	if video.Input.Height > video.Output.Video.Limits.VResMax || float64(video.Input.Width) > float64(video.Output.Video.Limits.VResMax) / 0.5625 {
		vertRes = video.Output.Video.Limits.VResMax
		scaleExpr := ""
		scaleAlgo := ""
		if float64(video.Input.Height) / float64(video.Input.Width) < 0.5625 { // 0.5625 = 16:9
			horizRes := int(float64(video.Output.Video.Limits.VResMax) / 1.125) * 2 // very hacky way of ensuring a multiple of 2
			scaleExpr = strconv.Itoa(horizRes) + ":-2"
		} else {
			scaleExpr = "-2:" + strconv.Itoa(video.Output.Video.Limits.VResMax)
		}
		if pass == 1 {
			scaleAlgo = "bilinear"
		} else {
			scaleAlgo = "spline"
		}
		filters = append(filters, "scale=" + scaleExpr + ":flags=" + scaleAlgo)
	}

	// HDR tonemapping
	if video.Input.IsHDR {
		if settings.Decoding.TonemapHWAccel {
			filters = append(filters, "format=p010,hwupload,tonemap_opencl=tonemap=mobius:format=p010,hwdownload,format=p010")
		} else {
			filters = append(filters, "zscale=transfer=linear,tonemap=mobius,zscale=transfer=bt709")
		}
	}

	// Pixfmt conversion
	if video.Input.Pixfmt != video.Output.Video.Encoder.Pixfmt {
		filters = append(filters, "format=" + video.Output.Video.Encoder.Pixfmt)
	}

	return strings.Join(filters, ","), vertRes, FPS
}