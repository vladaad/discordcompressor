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

func encFDKaac(outFilename string, video *settings.Video, input io.ReadCloser) {
	var options []string
	encoderSettings := strings.Split(video.Output.Audio.Encoder.Options, " ")

	// Encoding options
	if video.Output.Audio.Encoder.UsesBitrate {
		options = append(options, "-b", strconv.FormatFloat(video.Output.Audio.Bitrate, 'f', -1, 64))
	}
	if video.Output.Audio.Encoder.Options != "" {
		options = append(options, encoderSettings...)
	}
	// Input & output options
	options = append(options, "-")
	options = append(options, "-o", outFilename)

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Running
	if !settings.DryRun {
		cmd := exec.Command(settings.General.FDKaacExecutable, options...)

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
