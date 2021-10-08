package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func encFFmpeg(inFilename string, outFilename string, bitrate float64, audioTracks int, eOptions *settings.AudioEncoder, startingTime float64, totalTime float64) {
	var options []string
	encoderSettings := strings.Split(eOptions.Options, " ")

	input := decodeAudio(inFilename, audioTracks, startingTime, totalTime)

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
		"-c:a", eOptions.Encoder,
	)
	if eOptions.Options != "" {
		options = append(options, encoderSettings...)
	}
	if eOptions.UsesBitrate {
		options = append(options,
			"-b:a", strconv.FormatFloat(bitrate, 'f', -1, 64) + "k",
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