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
		"-show_entries", "stream=r_frame_rate:stream=height:stream=pix_fmt:format=duration:format=bit_rate",
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
		// Pix_fmt
		if strings.HasPrefix(line, "streams.stream.0.pix_fmt") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			stats.Pixfmt = splitString[len(splitString)-1]
		}
	}

	// Getting audio info
	aprobe := exec.Command(
		settings.General.FFprobeExecutable,
		"-loglevel", "quiet",
		"-of", "flat",
		"-select_streams", "a",
		"-show_entries", "stream=index:stream=codec_name:stream=bit_rate",
		filepath,
		)

	pipe, err = aprobe.StdoutPipe()
	if err != nil {
		panic("Failed to start stdout pipe")
	}

	err = aprobe.Start()
	if err != nil {
		panic("Failed to start FFprobe")
	}

	ascanner := bufio.NewScanner(pipe)
	totalStreams := 0

	for ascanner.Scan() {
		line := ascanner.Text()
		// Stream counter
		if strings.HasPrefix(line, "streams.stream." + strconv.Itoa(totalStreams)) {
			totalStreams += 1
		}
		// Codec name
		if strings.HasPrefix(line, "streams.stream.0.codec_name") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			stats.AudioCodec = splitString[len(splitString)-1]
		}
		// Audio bitrate (not in mkv :/)
		if strings.HasPrefix(line, "streams.stream.0.bit_rate") {
			cleanedLine := strings.ReplaceAll(line, "\"", "")
			splitString := strings.Split(cleanedLine, "=")
			if splitString[len(splitString)-1] == "N/A" {
				stats.AudioBitrate = -1
			} else {
				stats.AudioBitrate, _ = strconv.Atoi(splitString[len(splitString)-1])
			}
		}
	}
	stats.AudioTracks = totalStreams
	return stats
}