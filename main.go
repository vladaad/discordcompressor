package main

import (
	"flag"
	"github.com/google/uuid"
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/encoder"
	"github.com/vladaad/discordcompressor/encoder/audio"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var targetSize int
var targetStartingTime float64
var lastSeconds float64
var targetTotalTime float64
var customOutputFile string
var input string

var forceEncoder string
var forceAEncoder string
var forceContainer string

func init() {
	// Update checker
	utils.CheckForUpdates()
	// Settings dir creation
	err := os.MkdirAll(utils.SettingsDir(), 0755)
	if err != nil {
		panic("Failed to create settings directory")
	}
	// Log setup
	logFileName := utils.SettingsDir()
	logFileName += "/dcomp.log"
	file, err := os.Create(logFileName)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	// Version print
	log.Println("Starting DiscordCompressor version " + build.VERSION)
	// Check for FFmpeg and FFprobe
	checkForFF()

	// Parsing flags
	settingsFile := flag.String("settings", "", "Selects the settings file to be used")
	startTime := flag.Float64("ss", float64(0), "Sets the starting time")
	targetTime := flag.Float64("t", float64(0), "Sets the time to encode")
	lastXSeconds := flag.Float64("last", float64(0), "Sets the time from the end to encode")
	targetSizeMB := flag.Float64("size", float64(-1), "Sets the target size in MB")
	debug := flag.Bool("debug", false, "Prints extra info")
	customOutput := flag.String("o", "", "Outputs to a specific filename")
	forceEncode := flag.String("c:v", "", "Uses a specific encoder")
	forceAEncode := flag.String("c:a", "", "Uses a specific audio encoder")
	forceContaine := flag.String("f", "", "Uses a specific container") // don't mind the cut off letters thanks
	flag.Parse()
	// Settings loading
	// Forcing
	forceEncoder = *forceEncode
	forceAEncoder = *forceAEncode
	forceContainer = *forceContaine

	input = flag.Args()[0]
	targetStartingTime = *startTime
	targetTotalTime = *targetTime
	lastSeconds = *lastXSeconds
	customOutputFile = *customOutput
	settings.Debug = *debug

	// Load defaults
	if *targetSizeMB == float64(-1) {
		*targetSizeMB = settings.General.TargetSizeMB
	}
	targetSize = int(*targetSizeMB * 8388608) // 1024*1024*8 - in bits

	settings.LoadSettings(*settingsFile)

	if len(input) == 0 {
		log.Println("No input video specified, closing...")
		os.Exit(0)
	}
}

func main() {
	// Empty function for now, will be used for future features if needed
	compress(input)
}

func compress(inVideo string) bool {
	var wg sync.WaitGroup
	// Initialize variables
	video := initVideo()
	video.InFile = inVideo

	// Check if file exists
	if _, err := os.Stat(video.InFile); err != nil {
		panic(video.InFile + " doesn't exist")
	}

	// Video analysis, time calculation
	log.Println("Analyzing video...")
	video.Input = metadata.GetStats(video.InFile, false)
	video = metadata.CalculateTime(video, lastSeconds, targetStartingTime, targetTotalTime)

	// Bitrate calculation
	// min(target bps/time, max bps)
	video.Output.Bitrate.Total = int(math.Min(float64(targetSize)/video.Time.Duration, float64(settings.Encoding.MaxBitrate*1024)))

	// Encoder selecction
	video.Output.Force.Video, video.Output.Force.Audio, video.Output.Force.Container = forceEncoder, forceAEncoder, forceContainer
	video = metadata.SelectEncoder(video)
	video = metadata.CalcOverhead(video)

	// Audio bitrate calculation
	if video.Input.ATracks > 0 {
		video = metadata.CalcAudioBitrate(video)
		video.Output.AudioFile = audio.GenFilename(video)
		wg.Add(1)
		if settings.Encoding.FastMode {
			go audio.EncodeAudio(video, &wg)
		} else {
			log.Println("Encoding audio...")
			audio.EncodeAudio(video, &wg)
		}
	} else {
		video.Output.Bitrate.Video = video.Output.Bitrate.Total
	}

	if settings.Debug {
		log.Println("Video bitrate:", video.Output.Bitrate.Video/1024)
		log.Println("Audio bitrate:", video.Output.Bitrate.Audio/1024)
	}

	// Resolution analysis
	if !settings.Encoding.FastMode && settings.Encoding.AutoRes {
		log.Println("Automatically choosing resolution, this may take a while...")
		video = encoder.CalculateResolution(video)
		if settings.Debug {
			log.Println("Chosen vertical resolution:", video.Output.Settings.MaxVRes)
		}
	}

	// Output filename
	suffix := strings.ReplaceAll(settings.General.OutputSuffix, "%s", strconv.FormatFloat(float64(targetSize)/8388608, 'f', 1, 64))
	video.OutFile = strings.TrimSuffix(video.InFile, path.Ext(video.InFile)) + suffix + "." + video.Output.Settings.Container
	if customOutputFile != "" {
		video.OutFile = customOutputFile + "." + video.Output.Settings.Container
	}
	// Encoding
	log.Println("Encoding, pass 1/2")
	encoder.EncodeVideo(video, 1)
	wg.Wait()
	log.Println("Encoding, pass 2/2")
	encoder.EncodeVideo(video, 2)

	// Cleanup
	os.Remove(video.UUID + "-0.log")
	os.Remove(video.UUID + "-0.log.mbtree")
	os.Remove(video.Output.AudioFile)

	return true
}

func checkForFF() {
	exit := false
	check := []string{"ffmpeg", "ffprobe"}

	for i := range check {
		if !utils.CheckIfPresent(check[i]) {
			message := check[i] + " not installed or not added to PATH"
			if runtime.GOOS == "windows" {
				message = message + ", you can download it here: https://github.com/BtbN/FFmpeg-Builds/releases"
			}
			log.Println(message)
			exit = true
		}
	}

	if exit {
		os.Exit(1)
	}
}

func initVideo() *settings.Vid {
	// god fucking dammit
	video := new(settings.Vid)
	time := new(settings.Time)
	bitrates := new(settings.Bitrates)
	output := new(settings.Out)
	force := new(settings.Force)
	output.Bitrate = bitrates
	video.Output = output
	video.Output.Force = force
	video.Time = time
	video.UUID = uuid.New().String()

	return video
}
