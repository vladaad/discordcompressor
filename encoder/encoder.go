package encoder

import (
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

func Encode(filename string, pass int, videoBitrate int, audioMerge bool, audioBitrate int) bool {
	var options []string
	// Vars
	outputFilename := strings.TrimSuffix(filename, path.Ext(filename)) + " (compressed)." + settings.SelectedVEncoder.Container
	vEncoderOptions := strings.Split(settings.SelectedVEncoder.Options, " ")
	aEncoderOptions := strings.Split(settings.SelectedAEncoder.Options, " ")
	times := appendTimes()
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

	// If bitrate not set
	if audioMerge && pass != 1 {
		options = append(options, "-i", strings.TrimSuffix(filename, path.Ext(filename)) + ".audio." + settings.SelectedVEncoder.Container)
	}

	// Video filters
	filters := filters(pass)
	if settings.Original == false && filters != "" {
		options = append(options, "-vf", filters)
	}

	// Video encoding options
	options = append (options,
		"-c:v", settings.SelectedVEncoder.Encoder,
		settings.SelectedVEncoder.PresetCmd, settings.SelectedSettings.Preset,
		"-b:v", strconv.Itoa(videoBitrate) + "k",
	)
	options = append(options, vEncoderOptions...)
	options = append(options, "-g", strconv.FormatFloat(FPS * float64(settings.SelectedVEncoder.Keyint), 'f', -1, 64))

	// 2pass
	if pass != 0 {
		options = append(options, settings.SelectedVEncoder.PassCmd, strconv.Itoa(pass))
	}

	// Audio encoding options + mapping
	if pass == 1 {
		options = append(options,
			"-map", "0:v:0",
			"-an",
		)
	} else {
		if audioMerge {
			options = append(options,
				"-map", "1:a:0",
				"-map", "0:v:0",
				"-c:a", "copy",
			)
		} else {
			options = append(options,
				"-map", "0:v:0",
				"-map", "0:a:0",
				"-c:a", settings.SelectedAEncoder.Encoder,
				"-b:a", strconv.Itoa(audioBitrate) + "k",
			)
			if settings.SelectedAEncoder.Options != "" {
				options = append(options, aEncoderOptions...)
			}
		}
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

func appendTimes() []string {
	var times []string
	if settings.Starttime != float64(0) {
		times = append(times, "-ss", strconv.FormatFloat(settings.Starttime, 'f', -1, 64))
	}
	if settings.Time > float64(0) {
		times = append(times, "-t", strconv.FormatFloat(settings.Time, 'f', -1, 64))
	}
	return times
}

func filters(pass int) string {
	var filters []string
	var fpsfilters []string // scale,tmix,fps is faster than tmix,fps,scale
	var resfilters []string
	// FPS
	FPS = settings.VideoStats.FPS
	if float64(settings.SelectedLimits.FPSMax) < settings.VideoStats.FPS {
		if settings.Encoding.HalveDownFPS {
			for FPS > float64(settings.SelectedLimits.FPSMax) {
				FPS /= 2
			}
		} else {
			FPS = float64(settings.SelectedLimits.FPSMax)
		}
		if settings.Encoding.TmixDownFPS && pass != 1 {
			frames := settings.VideoStats.FPS / FPS
			if frames < 2 {frames = 2}
			fpsfilters = append(fpsfilters, "tmix=frames=" + strconv.FormatFloat(frames, 'f', 0, 64))
		}
		fpsfilters = append(fpsfilters, "fps=" + strconv.FormatFloat(FPS, 'f', -1, 64))
	}

	// Resolution
	if settings.SelectedLimits.VResMax < settings.VideoStats.Height {
		if pass == 1 {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax))
		} else {
			resfilters = append(resfilters, "scale=-2:" + strconv.Itoa(settings.SelectedLimits.VResMax) + ":flags=lanczos")
		}
	}

	if settings.Encoding.TmixDownFPS && pass != 1 {
		filters = append(filters, resfilters...)
		filters = append(filters, fpsfilters...)
	} else {
		filters = append(filters, fpsfilters...)
		filters = append(filters, resfilters...)
	}

	return strings.Join(filters, ",")
}