package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func encQaac(inFilename string, outFilename string, bitrate float64, audioTracks int, eOptions *settings.AudioEncoder, startingTime float64, totalTime float64) string {
	var options []string
	encoderSettings := strings.Split(eOptions.Options, " ")

	tempFilename := outFilename + "wav"
	extractAudio(inFilename, tempFilename, "", audioTracks, startingTime, totalTime)

	// Encoding options
	if eOptions.UsesBitrate {
		options = append(options, "-a", strconv.FormatFloat(bitrate, 'f', -1, 64))
	}
	if eOptions.Options != "" {
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
	return outFilename + "m4a"
}