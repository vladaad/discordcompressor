package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func decodeAudio (inFilename string, audioTracks int, startingTime float64, totalTime float64) io.ReadCloser {
	var options []string
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
	options = append(options, "-c:a", "pcm_s16le")
	// Trackmix
	if settings.MixTracks && audioTracks > 1 {
		var filter []string
		for i := 0; i < audioTracks; i++ {
			filter = append(filter, "[0:a:" + strconv.Itoa(i) + "]")
		}
		filter = append(filter, "amix=inputs=", strconv.Itoa(audioTracks))
		filter = append(filter, "[out]")
		options = append(options, "-filter_complex", strings.Join(filter, ""), "-map", "[out]")
	} else {
		options = append(options, "-map", "0:a:0")
	}

	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")
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
