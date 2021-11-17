package main

import (
	"flag"
	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/encoder/audio"
	"github.com/vladaad/discordcompressor/encoder/video"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

var reEncV bool
var reEncA bool
var targetSizeKbit float64
var targetStartingTime float64
var lastSeconds float64
var targetTotalTime float64
var input []string
var wg sync.WaitGroup
var runningInstances int

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

	// Parsing flags
	settingsFile := flag.String("settings", "", "Selects the settings file to be used")
	startTime := flag.Float64("ss", float64(0), "Sets the starting time")
	targetTime := flag.Float64("t", float64(0), "Sets the time to encode")
	lastXSeconds := flag.Float64("last", float64(0), "Sets the time from the end to encode")
	targetSize := flag.Float64("size", float64(-1), "Sets the target size in MB")
	debug := flag.Bool("debug", false, "Prints extra info")
	focus := flag.String("focus", "", "Sets the focus")
	original := flag.Bool("noscale", false, "Disables FPS limiting and scaling")
	mixTracks := flag.Bool("mixaudio", false, "Mixes all audio tracks into one")
	normalize := flag.Bool("normalize", false, "Normalizes audio volume")
	dryRun := flag.Bool("dryrun", false, "Just prints commands instead of running")
	reEncode := flag.String("reenc", "", "Re-encodes even when not needed. \"a\", \"v\" or \"av\"")
	flag.Parse()
	// Settings loading
	input = flag.Args()
	targetStartingTime = *startTime
	targetTotalTime = *targetTime
	lastSeconds = *lastXSeconds
	settings.Debug = *debug
	settings.Original = *original
	settings.Focus = *focus
	settings.DryRun = *dryRun

	// Reenc
	reEncA, reEncV = false, false
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
	if len(input) == 0 && !newSettings {
		utils.OpenURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
	}

	if len(input) == 0 {
		log.Println("No input video specified, closing...")
		os.Exit(0)
	}
	// load defaults of some settings
	if *targetSize == float64(-1) {
		*targetSize = settings.Encoding.SizeTargetMB
	}
	if *mixTracks {
		settings.Advanced.MixAudioTracks = *mixTracks
	}
	if *normalize {
		settings.Advanced.NormalizeAudio = *normalize
	}
	settings.TargetSize = *targetSize
	targetSizeKbit = *targetSize * 8192

	// enable batch mode - stdout
	if len(input) > 1 && settings.General.BatchModeThreads > 1 {settings.BatchMode = true}

	if settings.Debug || !settings.BatchMode {
		settings.ShowStdOut = true
	} else {
		settings.ShowStdOut = false
	}
}

func main() {
	for i := range input {
		// yes this is a mess
		for {
			if runningInstances < settings.General.BatchModeThreads {
				wg.Add(1)
				runningInstances += 1
				go compress(input[i])
				break
			}
			time.Sleep(time.Millisecond * 50)
		}
	}
	wg.Wait()
	if len(input) > 1 {log.Println("All files completed!")}
}

func compress(inVideo string) bool {
	var prefix string
	var totalTime float64
	var startingTime float64
	defer wg.Done()
	// Logging
	_, cleanName := path.Split(strings.ReplaceAll(inVideo, "\\", "/"))
	if settings.BatchMode{prefix = "[" + cleanName + "] "}

	log.Println("Compressing " + cleanName)

	// Generate UUID
	UUID := utils.GenUUID()

	// Video analysis
	log.Println(prefix + "Analyzing video...")
	videoStats := metadata.GetStats(inVideo, false)

	// Checking time
	totalTime, startingTime = targetTotalTime, targetStartingTime
	if settings.BatchMode && (totalTime != 0 || startingTime != 0) {
		log.Fatalln("Cannot use time arguments with batch mode except -last!")
	} else {
		if lastSeconds != 0 && (totalTime != 0 || startingTime != 0) {
			log.Println(prefix + "Cannot use -t or -ss with -last!")
			return false
		}
		// LastSeconds
		if lastSeconds != 0 {
			startingTime = videoStats.Duration - lastSeconds
			videoStats.Duration = lastSeconds
		} else { // ss+t
			if startingTime + totalTime > videoStats.Duration {
				log.Println(prefix + "Invalid length!")
				return false
			}
			if totalTime != 0 {
				videoStats.Duration = totalTime
			} else if startingTime != 0 {
				videoStats.Duration -= startingTime
			}
		}
	}

	if settings.Debug {
		log.Println("Input stats:")
		log.Println(strconv.Itoa(videoStats.Height) + "p " + strconv.FormatFloat(videoStats.FPS, 'f', -1, 64) + "fps")
		log.Println("Length: " + strconv.FormatFloat(videoStats.Duration, 'f', -1, 64) + " seconds")
		log.Println("Pixel format: " + videoStats.Pixfmt)
		log.Println("Audio tracks: " + strconv.Itoa(videoStats.AudioTracks))
		if videoStats.AudioTracks != 0 {
			log.Println(videoStats.AudioCodec + ", " + strconv.FormatFloat(videoStats.AudioBitrate, 'f', 1, 64) + "k")
			log.Println(strconv.Itoa(videoStats.SampleRate) + "hz " + strconv.Itoa(videoStats.AudioChannels) + "ch")
		}
	}

	// Total bitrate calc
	totalBitrateUncomp, err := metadata.CalcTotalBitrate(targetSizeKbit, videoStats.Duration)
	if err {
		return false
	}

	// Choosing target
	videoEncoder, audioEncoder, target, limits := metadata.SelectEncoder(totalBitrateUncomp)
	outTarget := new(video.OutTarget)

	// Overshoot compensation
	overhead := metadata.CalcOverhead(math.Min(float64(limits.FPSMax), videoStats.FPS), videoStats.Duration)
	if target.Encoder == "libx264" {
		overhead += metadata.CalcH264Overhead(videoStats.Duration)
	}
	totalBitrate := totalBitrateUncomp - overhead

	// AB calc & passthrough
	hasAudio := true
	outTarget.AudioBitrate = metadata.CalcAudioBitrate(totalBitrate, audioEncoder, videoStats.AudioChannels)
	outTarget.AudioPassthrough, outTarget.VideoPassthrough, outTarget.AudioBitrate = metadata.CheckStreamCompatibility(inVideo, outTarget.AudioBitrate, totalBitrate, videoStats, startingTime, totalTime, videoEncoder, audioEncoder)
	if reEncA {outTarget.AudioPassthrough = false}
	if reEncV {outTarget.VideoPassthrough = false}

	// Audio encoding
	var audioFile string
	if !outTarget.AudioPassthrough && videoStats.AudioTracks != 0 {
		log.Println(prefix + "Encoding audio...")
		outTarget.AudioBitrate, audioFile = audio.EncodeAudio(inVideo, UUID, outTarget.AudioBitrate, videoEncoder.Container, audioEncoder, videoStats, startingTime, totalTime)
		if audioFile == "" {
			outTarget.AudioBitrate = 0
			hasAudio = false
		}
	} else if !outTarget.AudioPassthrough {
		outTarget.AudioBitrate = 0
		hasAudio = false
	}

	// Video bitrate calc
	outTarget.VideoBitrate = totalBitrate - outTarget.AudioBitrate

	// Debug
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.FormatFloat(totalBitrateUncomp, 'f', 1, 64) + "k")
		log.Println("Adjusted target bitrate: " + strconv.FormatFloat(totalBitrate, 'f', 1, 64) + "k")
		log.Println("Overhead: " + strconv.FormatFloat(overhead, 'f', 1, 64) + "k")
		if videoStats.AudioTracks != 0 {
			log.Println("Calculated video bitrate: " + strconv.FormatFloat(outTarget.VideoBitrate, 'f', 1, 64) + "k")
			log.Println("Calculated audio bitrate: " + strconv.FormatFloat(outTarget.AudioBitrate, 'f', 1, 64) + "k")
		}
	}

	// Encode
	if videoEncoder.TwoPass && !outTarget.VideoPassthrough {
		log.Println(prefix + "Encoding, pass 1/2")
		video.Encode(inVideo, audioFile, UUID, 1, false, videoStats, videoEncoder, target, limits, outTarget, audioEncoder, startingTime, totalTime)
		log.Println(prefix + "Encoding, pass 2/2")
		video.Encode(inVideo, audioFile, UUID, 2, hasAudio, videoStats, videoEncoder, target, limits, outTarget, audioEncoder, startingTime, totalTime)
	} else {
		log.Println(prefix + "Encoding, pass 1/1")
		video.Encode(inVideo, audioFile, UUID,0, hasAudio, videoStats, videoEncoder, target, limits, outTarget, audioEncoder, startingTime, totalTime)
	}

	os.Remove(UUID + "-0.log")
	os.Remove(UUID + "-0.log.mbtree")

	if hasAudio{os.Remove(audioFile)}

	log.Println("Finished compressing " + cleanName + "!")

	runningInstances -= 1
	return true
}