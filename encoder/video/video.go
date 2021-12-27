package vidEnc

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type OutTarget struct {
	AudioPassthrough bool
	VideoPassthrough bool
	VideoBitrate     float64
	AudioBitrate     float64
}

func Encode(video *settings.Video, outFilename string, pass int) bool {
	var options []string
	// Vars
	vEncoderOptions := strings.Split(video.Output.Video.Encoder.Options, " ")
	times := metadata.AppendTimes(video)
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
	if video.Input.IsHDR && settings.Decoding.TonemapHWAccel {
		options = append(options, "-hwaccel", "opencl")
	}
	options = append(options, times...)
	options = append(options, "-i", video.Filename)

	// Audio append
	if video.Output.Audio.Bitrate > 0 && !video.Output.Audio.Passthrough {
		options = append(options, "-i", video.Output.Audio.Filename)
	}

	// Video encoding options
	metaVertRes, FPS := video.Input.Height, video.Input.FPS
	if !video.Output.Video.Passthrough {
		var filter string
		// Video filters
		filter, metaVertRes, FPS = filters(video, pass)
		if filter != "" {
			options = append(options, "-vf", filter)
		}
		options = append(options,
			"-c:v", video.Output.Video.Encoder.Encoder,
			video.Output.Video.Encoder.PresetCmd, video.Output.Video.Preset,
			"-b:v", strconv.FormatFloat(video.Output.Video.Bitrate, 'f', -1, 64)+"k",
			"-vsync", "vfr",
		)
		if video.Output.Video.Encoder.Options != "" {
			options = append(options, vEncoderOptions...)
		}
		options = append(options, "-g", strconv.FormatFloat(FPS*float64(video.Output.Video.Encoder.Keyint), 'f', 0, 64))
		// 2pass
		if pass != 0 {
			options = append(options, video.Output.Video.Encoder.PassCmd, strconv.Itoa(pass))
			options = append(options, "-passlogfile", video.UUID)
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
	if video.Output.Audio.Passthrough {
		options = append(options, "-map", "0:a:0")
	} else if video.Output.Audio.Bitrate > 0 {
		options = append(options, "-map", "1:a:0")
	} else {
		options = append(options, "-an")
	}
	if video.Output.Audio.Passthrough || video.Output.Audio.Bitrate > 0 {
		options = append(options, "-c:a", "copy")
	}
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")
	options = append(options, addMetadata(video, metaVertRes, FPS)...)

	// Faststart for MP4
	if strings.ToLower(video.Output.Video.Encoder.Container) == "mp4" {
		options = append(options, "-movflags", "+faststart")
	}

	// Don't output to file in 1st pass
	if pass != 1 {
		options = append(options, outFilename)
	} else {
		options = append(options, "-f", "null", utils.NullDir())
	}

	// WEBM H264+AAC workaround
	if pass != 1 && video.Output.Video.Encoder.Container == "webm" {
		options = append(options, "-f", "matroska")
	}

	if settings.Debug || settings.DryRun {
		log.Println(options)
	}

	// Execution
	if !settings.DryRun {
		cmd := exec.Command(settings.General.FFmpegExecutable, options...)

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

	return true
}
