package encoder

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func EncodeVideo(video *settings.Vid, pass int) {
	var options []string
	EncoderOptions := strings.Split(video.Output.Encoder.Args, " ")
	options = append(options, "-loglevel", "warning", "-stats")

	// Input
	options = append(options, metadata.AppendTimes(video)...)
	options = append(options, "-y", "-i")
	options = append(options, video.InFile)

	if video.Output.AudioFile != "" && pass == 2 {
		options = append(options, "-i", video.Output.AudioFile)
	}

	// Filtering
	var filters []string
	// Resolution
	if video.Output.Settings.MaxVRes < video.Input.Height {
		scaler := settings.Encoding.Scaler
		if pass == 1 {
			scaler = "fast_bilinear"
		}
		var filter string
		filter += "scale=-2:"
		filter += strconv.Itoa(video.Output.Settings.MaxVRes)
		filter += ":flags=" + scaler
		filters = append(filters, filter)
	}
	if filters != nil {
		combined := strings.Join(filters, ",")
		options = append(options, "-vf", combined)
	}
	// FPS
	video = calculateFPS(video)
	if video.Output.FPS != video.Input.FPS {
		options = append(options, fpsFilter(video)...)
		options = append(options, "-vsync", "cfr")
	} else {
		options = append(options, "-vsync", "vfr")
	}

	// Encoding
	options = append(options, EncoderOptions...)
	options = append(options, "-pix_fmt", video.Output.Encoder.Pixfmt)
	options = append(options, "-b:v", strconv.Itoa(video.Output.Bitrate.Video))

	options = append(options, "-pass", strconv.Itoa(pass))
	options = append(options, "-passlogfile", video.UUID)

	// Mapping
	options = append(options, "-map", "0:v:0")
	if video.Output.AudioFile != "" && pass == 2 {
		options = append(options, "-map", "1:a:0")
		options = append(options, "-c:a", "copy")
	}
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")

	// Faststart for MP4
	if strings.ToLower(video.Output.Settings.Container) == "mp4" {
		options = append(options, "-movflags", "+faststart")
	}

	// Output

	if pass != 1 {
		options = append(options, video.OutFile)
	} else {
		options = append(options, "-f", "null", "-")
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
}
