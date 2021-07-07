package main

import (
	"flag"
	"fmt"
	"github.com/vladaad/discordcompressor/encoder"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	// Log setup
	file, err := os.Create("dcomp.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Parsing flags
	settingsFile := flag.String("settings", "", "Selects the settings file to be used")
	inputVideo := flag.String("i", "", "Sets the input video")
	startTime := flag.Float64("ss", float64(0), "Sets the starting time")
	time := flag.Float64("t", float64(0), "Sets the time to encode")
	targetSize := flag.Float64("size", float64(-1), "Sets the target size in MB")
	debug := flag.Bool("debug", false, "Prints extra info")
	focus := flag.String("focus", "", "Sets the focus")
	original := flag.Bool("noscale", false, "Disables FPS limiting and scaling")
	mixTracks := flag.Bool("mixaudio", false, "Mixes all audio tracks into one")
	flag.Parse()

	// Settings loading
	settings.InputVideo = *inputVideo
	settings.Starttime = *startTime
	settings.Time = *time
	settings.Debug = *debug
	settings.Original = *original
	settings.Focus = *focus
	settings.MixTracks = *mixTracks

	// ;)
	newSettings := settings.LoadSettings(*settingsFile)
	if *inputVideo == "" && !newSettings {
		OpenURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	}

	if *inputVideo == "" {
		log.Println("No input video specified, closing...")
		os.Exit(0)
	}
	// targetSizeMB defaults loading
	if *targetSize == float64(-1) {
		*targetSize = settings.Encoding.SizeTargetMB
	}

	// Video analysis
	log.Println("Analyzing video...")
	settings.VideoStats = metadata.GetStats(*inputVideo, false)
	if settings.Debug {
		log.Println("Input stats:")
		log.Println(strconv.Itoa(settings.VideoStats.Height) + "p " + strconv.FormatFloat(settings.VideoStats.FPS, 'f', -1, 64) + "fps")
		log.Println("Length: " + strconv.FormatFloat(settings.VideoStats.Duration, 'f', -1, 64) + " seconds")
		log.Println("Pixel format: " + settings.VideoStats.Pixfmt)
		log.Println("Audio tracks: " + strconv.Itoa(settings.VideoStats.AudioTracks))
		if settings.VideoStats.AudioTracks != 0 {
			log.Println(settings.VideoStats.AudioCodec + ", " + strconv.FormatFloat(settings.VideoStats.AudioBitrate, 'f', 1, 64) + "k")
		}
	}
	// Choosing target, audio encoding
	settings.OutputTarget = metadata.CalculateTarget(*inputVideo, *targetSize * 8192)

	// Video bitrate calc
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.FormatFloat(settings.MaxTotalBitrate, 'f', 1, 64) + "k")
	}

	// Encode
	if settings.SelectedVEncoder.TwoPass && !settings.OutputTarget.VideoPassthrough {
		log.Println("Encoding, pass 1/2")
		encoder.Encode(*inputVideo, 1)
		log.Println("Encoding, pass 2/2")
		encoder.Encode(*inputVideo, 2)
	} else {
		log.Println("Encoding, pass 1/1")
		encoder.Encode(*inputVideo, 0)
	}
	log.Println("Cleaning up...")
	os.Remove("ffmpeg2pass-0.log")
	os.Remove("ffmpeg2pass-0.log.mbtree")
	os.Remove(strings.TrimSuffix(settings.InputVideo, path.Ext(settings.InputVideo)) + ".audio." + settings.SelectedVEncoder.Container)
	log.Println("Finished!")
}

func OpenURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}