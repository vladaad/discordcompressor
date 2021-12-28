package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func decodeAudio (video *settings.Video, lnParams *LoudnormParams) io.ReadCloser {
	var options []string
	dontDownmix := []int{1, 2, 6, 8}

	times := metadata.AppendTimes(video)
	if settings.Debug {
		options = append(options,
			"-loglevel", "warning", "-stats",
		)
	} else {
		options = append(options,
			"-loglevel", "quiet", "-stats",
		)
	}
	options = append(options, "-y")
	options = append(options, times...)
	options = append(options, "-i", video.Filename)

	// Encoding options
	options = append(options, "-c:a", "pcm_s32le")
	// Filters
	if video.Output.Audio.Mix && video.Input.AudioTracks < 2 {
		video.Output.Audio.Mix = false
	}
	filters, mapping := filters(video, lnParams)

	if filters != "" {options = append(options, "-filter_complex", filters)}
	options = append(options, "-map", mapping)

	// Mapping
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")

	if !utils.ContainsInt(video.Input.AudioChannels, dontDownmix) || (video.Output.Audio.Mix && video.Input.AudioTracks > 1) {
		options = append(options, "-ac", "2")
	}

	if strings.Contains(video.Output.Audio.Encoder.CodecName, "aac") {
		options = append(options, "-ar", "44100")
	} else {
		options = append(options, "-ar", "48000")
	}
	options = append(options, "-f", "wav", "-")

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Running
	cmd := exec.Command(settings.General.FFmpegExecutable, options...)
	pipe, _ := cmd.StdoutPipe()

	if settings.Debug {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	return pipe
}
