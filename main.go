package main

import (
	"flag"
	"github.com/kardianos/osext"
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/encoder/audio"
	"github.com/vladaad/discordcompressor/encoder/video"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	// Log setup
	logFileName, _ := osext.ExecutableFolder()
	logFileName += "/dcomp.log"
	file, err := os.Create(logFileName)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	// Version print
	log.Println("Starting discordcompressor version " + build.VERSION)

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
	inVideo := *inputVideo
	startingTime := *startTime
	totalTime := *time
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
	if inVideo == "" && !newSettings {
		utils.OpenURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	}

	if inVideo == "" {
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
	videoStats := metadata.GetStats(inVideo, false)
	// Checking time
	if startingTime + totalTime > videoStats.Duration {
		log.Println("Invalid length!")
		os.Exit(0)
	}
	if totalTime != 0 {
		videoStats.Duration = totalTime
	} else if startingTime != 0 {
		videoStats.Duration =- startingTime
	}
	if settings.Debug {
		log.Println("Input stats:")
		log.Println(strconv.Itoa(videoStats.Height) + "p " + strconv.FormatFloat(videoStats.FPS, 'f', -1, 64) + "fps")
		log.Println("Length: " + strconv.FormatFloat(videoStats.Duration, 'f', -1, 64) + " seconds")
		log.Println("Pixel format: " + videoStats.Pixfmt)
		log.Println("Audio tracks: " + strconv.Itoa(videoStats.AudioTracks))
		if videoStats.AudioTracks != 0 {
			log.Println(videoStats.AudioCodec + ", " + strconv.FormatFloat(videoStats.AudioBitrate, 'f', 1, 64) + "k")
		}
	}
	// Total bitrate calc
	totalBitrate := metadata.CalcTotalBitrate(targetSizeKbit, videoStats.Duration)
	// Choosing target
	videoEncoder, audioEncoder, target, limits := metadata.SelectEncoder(totalBitrate)
	log.Println(limits)
	outTarget := new(video.OutTarget)
	// AB calc & passthrough
	hasAudio := true
	outTarget.AudioBitrate = metadata.CalcAudioBitrate(totalBitrate, settings.AudioEncoder{})
	outTarget.AudioPassthrough, outTarget.VideoPassthrough, outTarget.AudioBitrate = metadata.CheckStreamCompatibility(inVideo, outTarget.AudioBitrate, totalBitrate, videoStats, startingTime, totalTime, videoEncoder, audioEncoder)
	if reEncA {outTarget.AudioPassthrough = false}
	if reEncV {outTarget.VideoPassthrough = false}
	// Audio encoding
	if !outTarget.AudioPassthrough && videoStats.AudioTracks != 0 {
		log.Println("Encoding audio...")
		outTarget.AudioBitrate, settings.AudioFile = audio.EncodeAudio(inVideo, outTarget.AudioBitrate, videoStats.AudioTracks, videoEncoder.Container, audioEncoder, startingTime, totalTime)
	} else if !outTarget.AudioPassthrough {
		outTarget.AudioBitrate = 0
		hasAudio = false
	}
	// Video bitrate calc
	outTarget.VideoBitrate = totalBitrate - outTarget.AudioBitrate
	// Debug
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.FormatFloat(totalBitrate, 'f', 1, 64) + "k")
		if videoStats.AudioTracks != 0 {
			log.Println("Calculated video bitrate: " + strconv.FormatFloat(outTarget.VideoBitrate, 'f', 1, 64) + "k")
			log.Println("Calculated audio bitrate: " + strconv.FormatFloat(outTarget.AudioBitrate, 'f', 1, 64) + "k")
		}
	}

	// Encode
	if videoEncoder.TwoPass && outTarget.VideoPassthrough {
		log.Println("Encoding, pass 1/2")
		video.Encode(inVideo, 1, false, videoStats, videoEncoder, target, limits, outTarget, startingTime, totalTime)
		log.Println("Encoding, pass 2/2")
		video.Encode(inVideo, 2, hasAudio, videoStats, videoEncoder, target, limits, outTarget, startingTime, totalTime)
	} else {
		log.Println("Encoding, pass 1/1")
		video.Encode(inVideo, 0, hasAudio, videoStats, videoEncoder, target, limits, outTarget, startingTime, totalTime)
	}
	log.Println("Cleaning up...")
	os.Remove("ffmpeg2pass-0.log")
	os.Remove("ffmpeg2pass-0.log.mbtree")
	os.Remove(settings.AudioFile)
	log.Println("Finished!")
}