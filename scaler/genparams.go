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
		pixfmt_cuda := pixfmt
		samePixfmt := false
		if targetPixfmt == "" {
			samePixfmt = true
			pixfmt = video.Input.Pixfmt
			pixfmt_cuda = pixfmt
		}
		if pixfmt_cuda == "yuv420p10le" {
			pixfmt_cuda = "p010"
		}

		filter += "hwupload_cuda,"
		filter += "scale_cuda=" + strconv.Itoa(targetWidth)
		filter += ":" + strconv.Itoa(targetHeight)
		if !samePixfmt {
			filter += ":format=" + pixfmt_cuda
		}
		filter += ",hwdownload,format=" + pixfmt_cuda //hwdownload always needs pixfmt

	}
	return filter
}
