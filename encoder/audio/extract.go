package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func extractAudio (inFilename string, outFilename string, encoder string) {
	var options []string
	times := metadata.AppendTimes()
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
	if encoder != "" {
		options = append(options,
			"-c:a", encoder,
		)
	}
	// Trackmix
	if settings.MixTracks && settings.VideoStats.AudioTracks > 1 {
		var filter []string
		for i := 0; i < settings.VideoStats.AudioTracks; i++ {
			filter = append(filter, "[0:a:" + strconv.Itoa(i) + "]")
		}
		filter = append(filter, "amix=inputs=", strconv.Itoa(settings.VideoStats.AudioTracks))
		filter = append(filter, "[out]")
		options = append(options, "-filter_complex", strings.Join(filter, ""), "-map", "[out]")
	} else {
		options = append(options, "-map", "0:a:0")
	}

	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")
	options = append(options, outFilename)

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Running
	if !settings.DryRun {
		cmd := exec.Command(settings.General.FFmpegExecutable, options...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			panic(err)
		}
	}
}
