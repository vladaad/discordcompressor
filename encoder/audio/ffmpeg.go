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

func encFFmpeg(inFilename string, outFilename string, bitrate float64, audioTracks int) {
	var options []string
	encoderSettings := strings.Split(settings.SelectedAEncoder.Options, " ")
	times := metadata.AppendTimes()

	tempFilename := inFilename + ".temp.wav"
	useTempFile := false
	if settings.MixTracks && audioTracks > 1 {
		extractAudio(inFilename, tempFilename, "", audioTracks)
		useTempFile = true
	}

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
	options = append(options, times...)
	if useTempFile {
		options = append(options, "-i", tempFilename)
	} else {
		options = append(options, "-i", inFilename)
	}

	// Encoding options
	options = append(options,
		"-c:a", settings.SelectedAEncoder.Encoder,
	)
	if settings.SelectedAEncoder.Options != "" {
		options = append(options, encoderSettings...)
	}
	if settings.SelectedAEncoder.UsesBitrate {
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

		if useTempFile {
			err = os.Remove(tempFilename)
			if err != nil {
				panic("Failed to remove temporary audio file")
			}
		}
	}
}