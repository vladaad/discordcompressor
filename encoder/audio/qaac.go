package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func encQaac(inFilename string, bitrate float64, audioTracks int) string {
	var options []string
	encoderSettings := strings.Split(settings.SelectedAEncoder.Options, " ")

	tempFilename := inFilename + ".temp.wav"
	extractAudio(inFilename, tempFilename, "", audioTracks)

	// Encoding options
	if settings.SelectedAEncoder.UsesBitrate {
		options = append(options, "-a", strconv.FormatFloat(bitrate, 'f', -1, 64))
	}
	if settings.SelectedAEncoder.Options != "" {
		options = append(options, encoderSettings...)
	}
	// Output options
	options = append(options, tempFilename)

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Running
	if !settings.DryRun {
		cmd := exec.Command(settings.General.QaacExecutable, options...)

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

		err = os.Remove(tempFilename)
		if err != nil {
				panic("Failed to remove temporary audio file")
		}
	}
	return inFilename + ".temp.m4a"
}