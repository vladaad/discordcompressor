package subtitles

import (
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"strconv"
	"strings"
)

func FindTime(filename string, find string, stream int) (start float64, end float64) {
	var args []string

	args = append(args, "-i", filename)
	args = append(args, "-map", "0:" + strconv.Itoa(stream))
	args = append(args, "-f", "ass")
	args = append(args, "-")

	subtitles := utils.CommandOutput(settings.General.FFmpegExecutable, args)
	lines := parseASS(subtitles)
	for i := range lines {
		if strings.Contains(lines[i].text, find) {
			return lines[i].startTime, lines[i].endTime
		}
	}

	return -1, -1
}