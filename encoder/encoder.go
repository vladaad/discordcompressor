package encoder

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var FPS float64

func Encode(filename string, pass int) bool {
	var options []string
	// Vars
	outputFilename := strings.TrimSuffix(filename, path.Ext(filename)) + " (compressed)." + settings.SelectedVEncoder.Container
	vEncoderOptions := strings.Split(settings.SelectedVEncoder.Options, " ")
	times := metadata.AppendTimes()
	// Command
	if settings.Debug {
		options = append(options,
			"-loglevel", "warning", "-stats",
		)
	} else {
		options = append(options,
			"-loglevel", "quiet", "-stats",
			)
	}
	options = append(options,
		"-y", "-hwaccel", settings.Decoding.HardwareAccel,
	)
	options = append(options, times...)
	options = append(options, "-i", filename)

	// Audio append
	if pass != 1 && settings.VideoStats.AudioTracks != 0 && !settings.OutputTarget.AudioPassthrough {
		options = append(options, "-i", strings.TrimSuffix(filename, path.Ext(filename)) + ".audio." + settings.SelectedVEncoder.Container)
	}

	// Video filters
	filters := filters(pass)
	if settings.Original == false && filters != "" {
		options = append(options, "-vf", filters)
	}

	// Video encoding options
	if !settings.OutputTarget.VideoPassthrough {
		options = append (options,
			"-c:v", settings.SelectedVEncoder.Encoder,
			settings.SelectedVEncoder.PresetCmd, settings.SelectedSettings.Preset,
			"-b:v", strconv.Itoa(settings.OutputTarget.VideoBitrate) + "k",
			"-vsync", "cfr",
		)
		if settings.SelectedVEncoder.Options != "" {
			options = append(options, vEncoderOptions...)
		}
		options = append(options, "-g", strconv.FormatFloat(FPS * float64(settings.SelectedVEncoder.Keyint), 'f', -1, 64))
		// 2pass
		if pass != 0 {
			options = append(options, settings.SelectedVEncoder.PassCmd, strconv.Itoa(pass))
		}
	} else {
		options = append(options,
			"-c:v", "copy",
		)
	}

	// Mapping
	options = append(options,
		"-map", "0:v:0",
	)
	if pass != 1 && settings.VideoStats.AudioTracks != 0 {
		if settings.OutputTarget.AudioPassthrough {
			options = append(options, "-map", "0:a:0")
		} else {
			options = append(options, "-map", "1:a:0")
		}
		options = append(options, "-c:a", "copy")
	} else {
		options = append(options, "-an")
	}

	// Faststart for MP4
	if strings.ToLower(settings.SelectedVEncoder.Container) == "mp4" {
		options = append(options, "-movflags", "+faststart")
	}

	// Don't output to file in 1st pass
	if pass != 1 {
		options = append(options, outputFilename)
	} else {
		var null string
		switch runtime.GOOS {
		case "windows":
			null = "NUL"
		default:
			null = "/dev/null"
		}
		options = append(options, "-f", "matroska", null) // -f null can break 2pass w/ mkv for whatever reason
	}

	if settings.Debug {
		log.Println(options)
	}

	// Execution
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

	return true
}