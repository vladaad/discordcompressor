package metadata

import (
	"strconv"
)

func AppendTimes(startingTime float64, totalTime float64) []string {
	var times []string
	if startingTime != float64(0) {
		times = append(times, "-ss", strconv.FormatFloat(startingTime, 'f', -1, 64))
	}
	if totalTime > float64(0) {
		times = append(times, "-t", strconv.FormatFloat(totalTime, 'f', -1, 64))
	}
	return times
}