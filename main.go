package main

import (
	"flag"
	"fmt"
	"github.com/vladaad/discordcompressor/encoder/audio"
	"github.com/vladaad/discordcompressor/encoder/video"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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
	dryRun := flag.Bool("dryrun", false, "Just prints commands instead of running")
	reEncode := flag.String("reenc", "", "Re-encodes even when not needed. \"a\", \"v\" or \"av\"")
	flag.Parse()

	// Settings loading
	settings.InputVideo = *inputVideo
	settings.Starttime = *startTime
	settings.Time = *time
	settings.Debug = *debug
	settings.Original = *original
	settings.Focus = *focus
	settings.MixTracks = *mixTracks
	settings.DryRun = *dryRun

	// Reenc
	reEncA, reEncV := false, false
	if !(*reEncode == "a" || *reEncode == "av" || *reEncode == "va" || *reEncode == "v" || *reEncode == "") {
		log.Println("The re-encode argument must be \"a\", \"v\" or \"av\"/\"va\".")
		os.Exit(0)
	} else {
		switch *reEncode {
		case "av", "va":
			reEncA, reEncV = true, true
		case "v":
			reEncV = true
		case "a":
			reEncA = true
		}
	}

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
	targetSizeKbit := *targetSize * 8192

	// Video analysis
	log.Println("Analyzing video...")
	settings.VideoStats = metadata.GetStats(*inputVideo, false)
	// Checking time
	if settings.Starttime + settings.Time > settings.VideoStats.Duration {
		log.Println("Invalid length!")
		os.Exit(0)
	}
	if settings.Time != 0 {
		settings.VideoStats.Duration = settings.Time
	} else if settings.Starttime != 0 {
		settings.VideoStats.Duration = settings.VideoStats.Duration - settings.Starttime
	}
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
	// Total bitrate calc
	settings.MaxTotalBitrate = metadata.CalcTotalBitrate(targetSizeKbit)
	// Choosing target
	metadata.SelectEncoder(settings.MaxTotalBitrate)
	t := new(settings.OutTarget)
	// AB calc & passthrough
	hasAudio := true
	t.AudioBitrate = metadata.CalcAudioBitrate(settings.MaxTotalBitrate)
	t.AudioPassthrough, t.VideoPassthrough, t.AudioBitrate = metadata.CheckStreamCompatibility(*inputVideo, t.AudioBitrate)
	if reEncA {t.AudioPassthrough = false}
	if reEncV {t.VideoPassthrough = false}
	// Audio encoding
	if !t.AudioPassthrough && settings.VideoStats.AudioTracks != 0 {
		log.Println("Encoding audio...")
		t.AudioBitrate, settings.AudioFile = audio.EncodeAudio(*inputVideo, t.AudioBitrate)
	} else if !t.AudioPassthrough {
		t.AudioBitrate = 0
		hasAudio = false
	}
	// Video bitrate calc
	t.VideoBitrate = settings.MaxTotalBitrate - t.AudioBitrate
	settings.OutputTarget = t
	// Debug
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.FormatFloat(settings.MaxTotalBitrate, 'f', 1, 64) + "k")
		if settings.VideoStats.AudioTracks != 0 {
			log.Println("Calculated video bitrate: " + strconv.FormatFloat(settings.OutputTarget.VideoBitrate, 'f', 1, 64) + "k")
			log.Println("Calculated audio bitrate: " + strconv.FormatFloat(settings.OutputTarget.AudioBitrate, 'f', 1, 64) + "k")
		}
	}

	// Encode
	if settings.SelectedVEncoder.TwoPass && !settings.OutputTarget.VideoPassthrough {
		log.Println("Encoding, pass 1/2")
		video.Encode(*inputVideo, 1, false)
		log.Println("Encoding, pass 2/2")
		video.Encode(*inputVideo, 2, hasAudio)
	} else {
		log.Println("Encoding, pass 1/1")
		video.Encode(*inputVideo, 0, hasAudio)
	}
	log.Println("Cleaning up...")
	os.Remove("ffmpeg2pass-0.log")
	os.Remove("ffmpeg2pass-0.log.mbtree")
	os.Remove(settings.AudioFile)
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