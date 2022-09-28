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

func encFHG(video *settings.Vid, input io.ReadCloser) {
	var options []string
	var encoderSettings []string
	if video.Output.AEncoder.Args != "" {
		encoderSettings = strings.Split(video.Output.AEncoder.Args, " ")
	}

	// encoding
	if encoderSettings != nil {
		options = append(options, encoderSettings...)
	}
	if !video.Output.AEncoder.TVBR {
		options = append(options, "--cbr", strconv.Itoa(video.Output.Bitrate.Audio/1024))
	}
	// input
	options = append(options, "-")
	// output
	options = append(options, video.Output.AudioFile)

	if settings.Debug {
		log.Println(options)
	}

	// running
	cmd := exec.Command(settings.General.FHGaacExecutable, options...)

	cmd.Stdin = input

	if !settings.Encoding.FastMode {
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
