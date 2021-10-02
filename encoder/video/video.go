package video

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

func Encode(filename string, pass int, audio bool, videoStats *metadata.VidStats, eOptions *settings.Encoder, eTarget *settings.Target, limit *settings.Limits, startingTime float64, totalTime float64) bool {
	var options []string
	// Vars
	outputFilename := strings.TrimSuffix(filename, path.Ext(filename)) + " (compressed)." + eOptions.Container
	vEncoderOptions := strings.Split(eOptions.Options, " ")
	times := metadata.AppendTimes(startingTime, totalTime)
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
	if audio && !settings.OutputTarget.AudioPassthrough {
		options = append(options, "-i", settings.AudioFile)
	}

	// Video filters
	filters := filters(pass, videoStats, limit, eOptions.Pixfmt)
	if filters != "" {
		options = append(options, "-vf", filters)
	}

	// Video encoding options
	if !settings.OutputTarget.VideoPassthrough {
		options = append(options,
			"-c:v", eOptions.Encoder,
			eOptions.PresetCmd, eTarget.Preset,
			"-b:v", strconv.FormatFloat(settings.OutputTarget.VideoBitrate, 'f', -1, 64)+"k",
			"-vsync", "vfr",
		)
		if eOptions.Options != "" {
			options = append(options, vEncoderOptions...)
		}
		options = append(options, "-g", strconv.FormatFloat(FPS*float64(eOptions.Keyint), 'f', 0, 64))
		// 2pass
		if pass != 0 {
			options = append(options, eOptions.PassCmd, strconv.Itoa(pass))
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
	if settings.OutputTarget.AudioPassthrough {
		options = append(options, "-map", "0:a:0")
	} else if audio {
		options = append(options, "-map", "1:a:0")
	} else {
		options = append(options, "-an")
	}
	if settings.OutputTarget.AudioPassthrough || audio {
		options = append(options, "-c:a", "copy")
	}
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")

	// Faststart for MP4
	if strings.ToLower(eOptions.Container) == "mp4" {
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

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Execution
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
	}

	return true
}
