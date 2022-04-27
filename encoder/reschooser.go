package encoder

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func CalculateResolution(video *settings.Vid) *settings.Vid {
	// options
	heights := []int{108, 144, 180, 240, 288, 360, 480, 540, 576, 648, 720, 900, 1080, 1152, 1440, 2160}
	var options []string
	var bdrate int
	testRes := math.Min(float64(video.Input.Height), 720)
	// bdrate selection depending on encoder
	switch utils.GetArg(video.Output.Encoder.Args, "-c:v") {
	case "libx264":
		switch utils.GetArg(video.Output.Encoder.Args, "-preset") {
		case "veryslow":
			bdrate = -27
		case "slower":
			bdrate = -25
		case "slow":
			bdrate = -22
		case "medium":
			bdrate = -20
		case "fast":
			bdrate = -15
		case "faster":
			bdrate = -10
		default:
			bdrate = -10
			log.Println("Your x264 preset was not optimized for automatic resolution selection")
			log.Println("Use faster to veryslow for best accuracy")
		}
	case "libvpx-vp9":
		args := []string{"-speed", "-cpu-used"}
		speed := 8
		for i := range args {
			spd, err := strconv.Atoi(utils.GetArg(video.Output.Encoder.Args, args[i]))
			if err == nil {
				speed = spd
				break
			}
		}
		if speed == 8 {
			log.Println("Your libvpx-vp9 preset was not detected for automatic resolution selection")
			log.Println("Assuming cpu-used 8")
		}
		bdrate = -55 + 3*speed
	case "libaom-av1":
		speed, err := strconv.Atoi(utils.GetArg(video.Output.Encoder.Args, "-cpu-used"))
		if err != nil {
			log.Println("Your aomenc preset was not detected for automatic resolution selection")
			log.Println("Assuming cpu-used 6")
			speed = 6
		}
		bdrate = -68 + speed
	default:
		bdrate = -10
	}
	bdratefac := 100 / (100 + float64(bdrate))
	tempFilename := video.UUID + ".resa.mkv"

	// fps
	video = calculateFPS(video)
	fpsf := fpsFilter(video)

	// ffmpeg options
	options = append(options, "-loglevel", "warning", "-stats")
	options = append(options, metadata.AppendTimes(video)...)
	options = append(options, "-i", video.InFile)
	options = append(options, "-map", "0:v:0")
	options = append(options, "-map_metadata", "-1")
	options = append(options, "-map_chapters", "-1")
	options = append(options, "-vf", "scale=-2:"+strconv.FormatFloat(testRes, 'f', 0, 64)+":flags=fast_bilinear")

	if fpsf != nil {
		options = append(options, fpsf...)
	}

	options = append(options, "-c:v", "libx264", "-preset", "veryfast", "-crf", "28")
	options = append(options, "-g", keyint(video))
	options = append(options, "-pix_fmt", "yuv420p", tempFilename)

	// encode
	cmd := exec.Command(settings.General.FFmpegExecutable, options...)
	if settings.Debug {
		log.Println(options)
	}
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

	obitrate := metadata.GetStats(tempFilename, false).Bitrate.Total

	// fancy calculation stolen from discordify.sh
	height := int(math.Sqrt((testRes * testRes) * bdratefac * (float64(video.Output.Bitrate.Video) * video.Time.Duration) / (float64(obitrate) * video.Time.Duration)))

	// rounding height
	minimumDiff := math.MaxInt
	finalHeight := video.Output.Settings.MaxVRes
	for i := range heights {
		if height >= heights[i] {
			diff := height - heights[i]
			if diff < minimumDiff {
				finalHeight = heights[i]
				minimumDiff = diff
			}
		} else {
			diff := heights[i] - height
			if diff < minimumDiff {
				finalHeight = heights[i]
				minimumDiff = diff
			}
		}
	}
	video.Output.Settings.MaxVRes = finalHeight

	os.Remove(tempFilename)

	return video
}
