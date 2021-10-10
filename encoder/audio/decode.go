package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"log"
	"os"
	"os/exec"
)

func decodeAudio (inFilename string, startingTime float64, totalTime float64, videoStats *metadata.VidStats) io.ReadCloser {
	var options []string
	dontDownmix := []int{1, 2, 6, 8}

	times := metadata.AppendTimes(startingTime, totalTime)
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
	options = append(options, "-i", inFilename)

	// Encoding options
	options = append(options, "-c:a", "pcm_s32le")
	// Filters
	mixTracks := false
	if settings.MixTracks && videoStats.AudioTracks > 1 {
		mixTracks = true
	}
	filters, mapping := filters(mixTracks, videoStats)

	if filters != "" {options = append(options, "-filter_complex", filters)}
	options = append(options, "-map", mapping)

	// Mapping
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")

	if !utils.ContainsInt(videoStats.AudioChannels, dontDownmix) || mixTracks {
		options = append(options, "-ac", "2")
	}
	options = append(options, "-ar", "48000")
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
