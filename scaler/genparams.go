package scaler

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func GenerateFilter(video *settings.Vid, targetWidth int, targetHeight int, targetPixfmt string, scaler string) string {
	var filter string
	if scaler != "cuda" { //swscale
		filter += "scale=" + strconv.Itoa(targetWidth)
		filter += ":" + strconv.Itoa(targetHeight)
		filter += ":flags=" + scaler
		if targetPixfmt != "" {
			filter += ",format=" + targetPixfmt
		}
	} else { //cuda
		pixfmt := targetPixfmt
		samePixfmt := false
		if targetPixfmt == "" {
			samePixfmt = true
			pixfmt = video.Input.Pixfmt
		}

		filter += "hwupload_cuda,"
		filter += "scale_cuda=" + strconv.Itoa(targetWidth)
		filter += ":" + strconv.Itoa(targetHeight)
		if !samePixfmt {
			filter += ":format=" + pixfmt
		}
		filter += ",hwdownload,format=" + pixfmt //hwdownload always needs pixfmt

	}
	return filter
}
