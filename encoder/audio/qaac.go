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

func encQaac(outFilename string, bitrate float64, eOptions *settings.AudioEncoder, input io.ReadCloser) {
	var options []string
	encoderSettings := strings.Split(eOptions.Options, " ")

	// Encoding options
	if eOptions.UsesBitrate {
		options = append(options, "-a", strconv.FormatFloat(bitrate, 'f', -1, 64))
	}
	if eOptions.Options != "" {
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
		cmd := exec.Command(settings.General.QaacExecutable, options...)

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