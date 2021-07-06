package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func EncodeAudio(filename string, bitrate float64) float64 {
	var options []string
	outputFilename := strings.TrimSuffix(filename, path.Ext(filename)) + ".audio." + settings.SelectedVEncoder.Container
	encoderSettings := strings.Split(settings.SelectedAEncoder.Options, " ")
	times := AppendTimes()

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
	options = append(options, "-i", filename)

	// Encoding options
	options = append(options, "-map", "0:a:0")
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

	// Output
	options = append(options, outputFilename)

	if settings.Debug {
		log.Println(options)
	}

	// Running

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

	return GetStats(outputFilename, true).Bitrate
}
