package metadata

import (
	"bufio"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"os/exec"
	"strconv"
	"strings"
)

func AnalyzeAudio(filename string) float64 {
	extract := exec.Command(
		settings.General.FFmpegExecutable,
		"-y", "-i", filename,
		"-map", "0:a:0",
		"-c:a", "copy",
		"-f", "nut",
		utils.NullDir(),
	)

	pipe, err := extract.StderrPipe()
	scanner := bufio.NewScanner(pipe)
	if err != nil {
		panic(err)
	}

	err = extract.Start()
	if err != nil {
		panic("Couldn't start FFmpeg")
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "bitrate=") {
			split := strings.Split(line, "bitrate=")
			split2 := strings.Split(split[len(split)-1], "/")
			replaced := strings.ReplaceAll(split2[0], "kbits", "")
			parsed, err := strconv.ParseFloat(strings.TrimSpace(replaced), 64)
			if err != nil {
				return 9999
			}
			return parsed
		}
	}
	return 9999
}
