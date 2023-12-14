package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os/exec"
)

func decodeAudio(video *settings.Vid) io.ReadCloser {
	var options []string

	// input
	options = append(options, "-loglevel", "warning", "-stats")
	options = append(options, metadata.AppendTimes(video)...)
	options = append(options, "-i", video.InFile)

	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")

	// filters
	filters, mapping := getFilters(video)
	if filters != "" {
		options = append(options, "-filter_complex", filters)
	}
	options = append(options, "-map", mapping)

	// output format
	options = append(options, "-c:a", "pcm_s32le")
	if video.Input.AChannels > 2 {
		options = append(options, "-ac", "2")
	}

	options = append(options, "-f", "wav", "-")

	if settings.Debug {
		log.Println(options)
	}

	cmd := exec.Command(settings.General.FFmpegExecutable, options...)
	pipe, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	return pipe
}
