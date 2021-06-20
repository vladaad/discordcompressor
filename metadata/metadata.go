package metadata

import (
	"bufio"
	"github.com/vladaad/discordcompressor/settings"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetStats(filepath string) *settings.VidStats {
	if _, err := os.Stat(filepath); err != nil {
		panic(filepath + " doesn't exist")
	}

	stats := new(settings.VidStats)

	probe := exec.Command(
		settings.General.FFprobeExecutable,
		"-loglevel", "quiet",
		"-of", "flat",
		"-select_streams", "v:0",
		"-show_entries", "stream=r_frame_rate:stream=height:format=duration:format=bit_rate",
		filepath,
		)

	pipe, err := probe.StdoutPipe()
	if err != nil {
		panic("Failed to start stdout pipe")
	}

	scanner := bufio.NewScanner(pipe)

	err = probe.Start()
	if err != nil {
		panic("Failed to start FFprobe")
	}

	for scanner.Scan() {

		line := scanner.Text()
		// Height
		if strings.HasPrefix(line, "streams.stream.0.height=") {
			splitString := strings.Split(line, "=")
			stats.Height, _ = strconv.Atoi(splitString[len(splitString)-1])
		}
		// FPS
		if strings.HasPrefix(line, "streams.stream.0.r_frame_rate") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			splitString = strings.Split(splitString[len(splitString)-1], "/")
			numerator, _ := strconv.ParseFloat(splitString[len(splitString)-2], 64)
			denominator, _ := strconv.ParseFloat(splitString[len(splitString)-1], 64)
			dividedFPS := numerator / denominator
			stats.FPS = dividedFPS
		}
		// Total bitrate
		if strings.HasPrefix(line, "format.bit_rate") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			stats.Bitrate, _ = strconv.Atoi(splitString[len(splitString)-1])
		}
		// Duration
		if strings.HasPrefix(line, "format.duration") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			stats.Duration, _ = strconv.ParseFloat(splitString[len(splitString)-1], 64)
		}
	}

	return stats
}