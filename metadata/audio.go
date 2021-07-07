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
	// Track mixing
	if settings.MixTracks && settings.VideoStats.AudioTracks > 1 {
		var filter []string
		for i := 0; i < settings.VideoStats.AudioTracks; i++ {
			filter = append(filter, "[0:a:" + strconv.Itoa(i) + "]")
		}
		filter = append(filter, "amix=inputs=", strconv.Itoa(settings.VideoStats.AudioTracks))
		filter = append(filter, "[out]")
		options = append(options, "-filter_complex", strings.Join(filter, ""), "-map", "[out]")
	} else {
		options = append(options, "-map", "0:a:0")
	}
	options = append(options, "-map_metadata", "-1")

	// Output
	options = append(options, outputFilename)

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
	}
	if settings.DryRun {
		if bitrate != 0 {
			return bitrate
		} else {
			return settings.SelectedAEncoder.MaxBitrate
		}
	} else {
		return GetStats(outputFilename, true).Bitrate
	}
}
