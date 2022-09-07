package main

import (
	"flag"
	"io"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vladaad/discordcompressor/build"
	"github.com/vladaad/discordcompressor/encoder/audio"
	vidEnc "github.com/vladaad/discordcompressor/encoder/video"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/subtitles"
	"github.com/vladaad/discordcompressor/utils"
)

var reEncV bool
var reEncA bool
var targetSizeKbit float64
var targetStartingTime float64
var lastSeconds float64
var targetTotalTime float64
var stringToFind string
var customOutputFile string
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
	// Check for FFmpeg and FFprobe
	checkForFF()

	// Parsing flags
	settingsFile := flag.String("settings", "", "Selects the settings file to be used")
	startTime := flag.Float64("ss", float64(0), "Sets the starting time")
	targetTime := flag.Float64("t", float64(0), "Sets the time to encode")
	lastXSeconds := flag.Float64("last", float64(0), "Sets the time from the end to encode")
	stringToFindA := flag.String("subfind", "", "Finds and cuts out time from subtitle text")
	targetSize := flag.Float64("size", float64(-1), "Sets the target size in MB")
	debug := flag.Bool("debug", false, "Prints extra info")
	focus := flag.String("focus", "", "Sets the focus")
	original := flag.Bool("noscale", false, "Disables FPS limiting and scaling")
	mixTracks := flag.Bool("mixaudio", false, "Mixes all audio tracks into one")
	normalize := flag.Bool("normalize", false, "Normalizes audio volume")
	dryRun := flag.Bool("dryrun", false, "Just prints commands instead of running")
	reEncode := flag.String("reenc", "", "Re-encodes even when not needed. \"a\", \"v\" or \"av\"")
	forceBenchScore := flag.Float64("forcescore", -1, "Forces a specific benchmark score when generating settings")
	customOutput := flag.String("o", "", "Outputs to a specific filename")
	flag.Parse()
	// Settings loading
	input = flag.Args()
	targetStartingTime = *startTime
	targetTotalTime = *targetTime
	lastSeconds = *lastXSeconds
	stringToFind = *stringToFindA
	customOutputFile = *customOutput
	settings.ForceScore = *forceBenchScore
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

	// batch mode checks
	if customOutputFile != "" && len(input) > 1 {
		log.Println("Can't output to the same file multiple times!")
		os.Exit(0)
	}

	// enable batch mode - stdout
	if len(input) > 1 && settings.General.BatchModeThreads > 1 {
		settings.BatchMode = true
	}

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
	if len(input) > 1 {
		log.Println("All files completed!")
	}
}

func compress(inVideo string) bool {
	var prefix string
	defer wg.Done()
	// Logging
	_, cleanName := path.Split(strings.ReplaceAll(inVideo, "\\", "/"))
	if settings.BatchMode {
		prefix = "[" + cleanName + "] "
	}
	log.Println("Compressing " + cleanName)

	// Initialize video
	video := initVideo()
	video.Filename = inVideo
	video.Size = targetSizeKbit

	video.Output.Audio.Normalize = settings.Advanced.NormalizeAudio
	video.Output.Audio.Mix = settings.Advanced.MixAudioTracks

	// Video analysis
	log.Println(prefix + "Analyzing video...")
	video.Input = metadata.GetStats(inVideo, false)

	// Subtitle checking
	if stringToFind != "" {
		if video.Input.MatchingSubs {
			targetStartingTime, targetTotalTime = subtitles.FindTime(video, stringToFind)
			if targetStartingTime == -1 || targetTotalTime == -1 {
				log.Println("Segment not found, try searching again! Keep in mind that discordcompressor can only find one specific subtitle for now.")
				os.Exit(0)
			}
			// Time compensation
			targetStartingTime -= settings.Advanced.SubStartOffset
			targetTotalTime += settings.Advanced.SubStartOffset + settings.Advanced.SubEndOffset
			// Clamping values
			targetTotalTime = math.Min(targetTotalTime, video.Input.Duration)
			targetStartingTime = math.Max(0, targetStartingTime)

			targetTotalTime -= targetStartingTime

		} else {
			log.Println("Error: subtitles with your target language not found.")
			os.Exit(0)
		}
	}

	// Checking time
	video.Time.Time, video.Time.Start = targetTotalTime, targetStartingTime
	if targetTotalTime == 0 {
		video.Time.Time = video.Input.Duration
	}
	if settings.BatchMode && (video.Time.Time != video.Input.Duration || video.Time.Start != 0) {
		log.Fatalln("Cannot use time arguments with batch mode except -last!")
	} else {
		if lastSeconds != 0 && (video.Time.Time != video.Input.Duration || video.Time.Start != 0) {
			log.Println(prefix + "Cannot use -t or -ss with -last!")
			return false
		}
		// LastSeconds
		if lastSeconds != 0 {
			video.Time.Start = video.Input.Duration - lastSeconds
			video.Time.Time = lastSeconds
		} else { // ss+t
			if video.Time.Start != 0 && video.Time.Time == video.Input.Duration {
				video.Time.Time = video.Input.Duration - video.Time.Start
			}
			if video.Time.Start+video.Time.Time > video.Input.Duration {
				log.Println(prefix + "Invalid length!")
				return false
			}
		}
	}

	if settings.Debug {
		log.Println("Input stats:")
		log.Println(strconv.Itoa(video.Input.Height) + "p " + strconv.FormatFloat(video.Input.FPS, 'f', -1, 64) + "fps")
		log.Println("Length: " + strconv.FormatFloat(video.Input.Duration, 'f', -1, 64) + " seconds")
		log.Println("Pixel format: " + video.Input.Pixfmt)
		log.Println("Audio tracks: " + strconv.Itoa(video.Input.AudioTracks))
		if video.Input.AudioTracks != 0 {
			log.Println(video.Input.AudioCodec + ", " + strconv.FormatFloat(video.Input.AudioBitrate, 'f', 1, 64) + "k")
			log.Println(strconv.Itoa(video.Input.SampleRate) + "hz " + strconv.Itoa(video.Input.AudioChannels) + "ch")
		}
	}

	// Total bitrate calc
	err := false
	video.Output.TotalBitrate, err = metadata.CalcTotalBitrate(video)
	if err {
		return false
	}

	// Choosing target
	video = metadata.SelectEncoder(video)

	// Overshoot compensation
	overhead := metadata.CalcOverhead(math.Min(float64(video.Output.Video.Limits.FPSMax), video.Input.FPS), video.Time.Time)
	if video.Output.Video.Encoder.Name == "libx264" {
		overhead += metadata.CalcH264Overhead(video.Time.Time)
	}
	video.Output.TotalBitrate = video.Output.TotalBitrate - overhead

	// AB calc & passthrough
	hasAudio := true
	video.Output.Audio.Bitrate = metadata.CalcAudioBitrate(video)
	video = metadata.CheckStreamCompatibility(video)
	if reEncA {
		video.Output.Audio.Passthrough = false
	}
	if reEncV {
		video.Output.Video.Passthrough = false
	}

	// Audio encoding
	if !video.Output.Audio.Passthrough && video.Input.AudioTracks != 0 {
		log.Println(prefix + "Encoding audio...")
		video.Output.Audio.Bitrate, video.Output.Audio.Filename = audio.EncodeAudio(video)
		if video.Output.Audio.Filename == "" {
			video.Output.Audio.Bitrate = 0
			hasAudio = false
		}
	} else if !video.Output.Audio.Passthrough {
		video.Output.Audio.Bitrate = 0
	}

	// Video bitrate calc
	video.Output.Video.Bitrate = video.Output.TotalBitrate - video.Output.Audio.Bitrate

	// Debug
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.FormatFloat(video.Output.TotalBitrate, 'f', 1, 64) + "k")
		log.Println("Overhead: " + strconv.FormatFloat(overhead, 'f', 1, 64) + "k")
		if video.Input.AudioTracks != 0 {
			log.Println("Calculated video bitrate: " + strconv.FormatFloat(video.Output.Video.Bitrate, 'f', 1, 64) + "k")
			log.Println("Calculated audio bitrate: " + strconv.FormatFloat(video.Output.Audio.Bitrate, 'f', 1, 64) + "k")
		}
	}

	suffix := strings.ReplaceAll(settings.General.OutputSuffix, "%s", strconv.FormatFloat(video.Size/8192, 'f', -1, 64))
	outFilename := strings.TrimSuffix(video.Filename, path.Ext(video.Filename)) + suffix + "." + video.Output.Video.Encoder.Container

	// Custom output filename
	if customOutputFile != "" {
		outFilename = customOutputFile + "." + video.Output.Video.Encoder.Container
	}

	// Subtitle extraction
	video.Output.Subs.BurnSubs = settings.Advanced.BurnSubtitles
	if !video.Input.MatchingSubs {
		video.Output.Subs.BurnSubs = false
	}
	if video.Output.Subs.BurnSubs {
		log.Println("Extracting subtitles...")
		video.Output.Subs.SubFile = subtitles.ExtractSubs(video)
		if video.Output.Subs.SubFile != "" {
			video.Input.SubtitleStream = metadata.GetStats(video.Output.Subs.SubFile, true).SubtitleStream // hacky but works
		} else {
			log.Println("Subtitles couldn't be extracted! Not burning")
		}
	}

	// Software HDR warning
	if video.Input.IsHDR && !settings.Decoding.TonemapHWAccel {
		log.Println(prefix + "Warning: tonemapping HDR video in software - this is very slow")
	}

	// Encode
	if video.Output.Video.Encoder.TwoPass && !video.Output.Video.Passthrough {
		log.Println(prefix + "Encoding, pass 1/2")
		vidEnc.Encode(video, outFilename, 1)
		log.Println(prefix + "Encoding, pass 2/2")
		vidEnc.Encode(video, outFilename, 2)
	} else {
		log.Println(prefix + "Encoding, pass 1/1")
		vidEnc.Encode(video, outFilename, 0)
	}

	os.Remove(video.Output.Subs.SubFile)
	os.Remove(video.UUID + "-0.log")
	os.Remove(video.UUID + "-0.log.mbtree")
	os.Remove("x264_lookahead.clbin")

	if hasAudio {
		os.Remove(video.Output.Audio.Filename)
	}

	log.Println("Finished compressing " + cleanName + "!")

	runningInstances -= 1
	return true
}

func checkForFF() {
	exit := false
	check := []string{"ffmpeg", "ffprobe"}

	for i := range check {
		if !utils.CommandExists(check[i]) {
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

func initVideo() *settings.Video {
	// god fucking dammit
	video := new(settings.Video)
	time := new(settings.Time)
	output := new(settings.Out)
	videoo := new(settings.VideoOut)
	audioo := new(settings.AudioOut)
	subs := new(settings.SubOut)
	output.Video, output.Audio, output.Subs = videoo, audioo, subs
	video.Time = time
	video.Output = output
	video.UUID = utils.GenUUID()

	return video
}
