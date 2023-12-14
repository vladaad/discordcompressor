package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func AppendTimes(video *settings.Vid) []string {
	var times []string
	if video.Time.Start != float64(0) {
		times = append(times, "-ss", strconv.FormatFloat(video.Time.Start, 'f', -1, 64))
	}
	times = append(times, "-t", strconv.FormatFloat(video.Time.Duration, 'f', -1, 64))
	return times
}
