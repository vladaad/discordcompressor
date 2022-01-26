package subtitles

import (
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"strconv"
	"strings"
)

func FindTime(video *settings.Video, find string) (start float64, end float64) {
	var args []string

	args = append(args, "-i", video.Filename)
	args = append(args, "-map", "0:"+strconv.Itoa(video.Input.SubtitleStream))
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
