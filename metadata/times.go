package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func AppendTimes() []string {
	var times []string
	if settings.Starttime != float64(0) {
		times = append(times, "-ss", strconv.FormatFloat(settings.Starttime, 'f', -1, 64))
	}
	if settings.Time > float64(0) {
		times = append(times, "-t", strconv.FormatFloat(settings.Time, 'f', -1, 64))
	}
	return times
}