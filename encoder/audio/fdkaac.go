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

func encFDK(video *settings.Vid, input io.ReadCloser) {
	var options []string
	encoderSettings := strings.Split(video.Output.AEncoder.Args, " ")
	// input
	options = append(options, "-")
	// encoding
	options = append(options, encoderSettings...)
	if !video.Output.AEncoder.TVBR {
		options = append(options, "-b", strconv.Itoa(video.Output.Bitrate.Audio))
	}
	// output
	options = append(options, "-o", video.Output.AudioFile)

	if settings.Debug {
		log.Println(options)
	}

	// running
	cmd := exec.Command(settings.General.FDKaacExecutable, options...)

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
