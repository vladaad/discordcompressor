package metadata

import (
	"bufio"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func PassthroughCheck(video *settings.Vid) *settings.Vid {
	webmCodecs := []string{"opus", "vorbis", "av1", "h264", "vp9"} // probably could've done this cleaner
	// Audio passthrough
	if video.Input.Bitrate.Audio == 0 && settings.Encoding.ForceGetABR {
		video.Input.Bitrate.Audio = AnalyzeAudio(video.InFile)
		video.Input.Bitrate.Video = video.Input.Bitrate.Total - video.Input.Bitrate.Audio
	}

	if video.Input.Bitrate.Audio != 0 && video.Input.Bitrate.Audio < video.Output.Bitrate.Audio {

		if video.Input.ACodec == "aac" && video.Output.Settings.Container == "mp4" {
			video.Output.APassthrough = true
		}

		if contains(video.Input.ACodec, webmCodecs) {
			if video.Output.Settings.Container == "mp4" && utils.GetArg(video.Output.Encoder.Args, "-c:v") == "libx264" {
				video.Output.Settings.Container = "webm"
				video.Output.APassthrough = true
				log.Println("Setting container to WebM for audio passthrough")
			} else if video.Output.Settings.Container == "webm" {
				video.Output.APassthrough = true
			}
		}
	}
	// Correct A/V bitrate if passthrough
	if video.Output.APassthrough {
		video.Output.Bitrate.Audio = video.Input.Bitrate.Audio
		video.Output.Bitrate.Video = video.Output.Bitrate.Total - video.Output.Bitrate.Audio
	}
	return video
}

func contains(str string, list []string) bool {
	for i := range list {
		if str == list[i] {
			return true
		}
	}
	return false
}

func AnalyzeAudio(filename string) int {
	extract := exec.Command(
		settings.General.FFmpegExecutable,
		"-y", "-i", filename,
		"-map", "0:a:0",
		"-c:a", "copy",
		"-f", "nut",
		"-",
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
				return 0
			}
			return int(parsed * 1024)
		}
	}
	return 0
}
