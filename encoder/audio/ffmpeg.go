package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func encFFmpeg(outFilename string, video *settings.Video, input io.ReadCloser) {
	var options []string
	encoderSettings := strings.Split(video.Output.Audio.Encoder.Options, " ")

	// Input options
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
	options = append(options, "-i", "-")

	// Encoding options
	options = append(options,
		"-c:a", video.Output.Audio.Encoder.Encoder,
	)
	if video.Output.Audio.Encoder.Options != "" {
		options = append(options, encoderSettings...)
	}
	if video.Output.Audio.Encoder.UsesBitrate {
		options = append(options,
			"-b:a", strconv.FormatFloat(video.Output.Audio.Bitrate, 'f', -1, 64)+"k",
		)
	}
	options = append(options, "-map", "0:a:0")

	// Output
	options = append(options, outFilename)

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Running
	if !settings.DryRun {
		cmd := exec.Command(settings.General.FFmpegExecutable, options...)

		cmd.Stdin = input

		if settings.ShowStdOut {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

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
