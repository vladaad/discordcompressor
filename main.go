package main

import (
	"flag"
	"github.com/vladaad/discordcompressor/encoder"
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var audioBitrate int
var videoBitrate int
var audioMerge bool

func main() {
	// Log setup
	file, err := os.Create("dcomp.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	settings.LoadSettings("")

	// Parsing flags
	inputVideo := flag.String("i", "", "Sets the input video")
	startTime := flag.Float64("ss", float64(0), "Sets the starting time")
	time := flag.Float64("t", float64(0), "Sets the time to encode")
	targetSizeMB := flag.Float64("size", settings.Encoding.SizeTargetMB, "Sets the target size in MB")
	debug := flag.Bool("debug", false, "Prints extra info")
	focus := flag.String("focus", "", "Sets the focus")
	original := flag.Bool("noscale", false, "Disables FPS limiting and scaling")
	flag.Parse()

	settings.InputVideo = *inputVideo
	settings.Starttime = *startTime
	settings.Time = *time
	settings.Debug = *debug
	settings.Original = *original
	settings.Focus = *focus

	targetSizeKbit := int(*targetSizeMB * float64(8192))

	// Video analysis
	log.Println("Analyzing video...")
	settings.VideoStats = metadata.GetStats(*inputVideo)
	if settings.Debug {
		log.Println("Input stats:")
		log.Println(strconv.Itoa(settings.VideoStats.Height) + "p " + strconv.FormatFloat(settings.VideoStats.FPS, 'f', -1, 64) + "fps")
		log.Println("Length: " + strconv.FormatFloat(settings.VideoStats.Duration, 'f', -1, 64) + " seconds")
	}

	// ss+t fixing
	if settings.Starttime + settings.Time > settings.VideoStats.Duration {
		log.Println("Invalid start or end time arguments!")
		os.Exit(0)
	}

	// Total bitrate calc
	if settings.Time != float64(0) {
		settings.VideoStats.Duration = settings.Time
	}
	totalBitrate := CalculateBitrate(settings.VideoStats, targetSizeKbit)
	if totalBitrate > settings.Encoding.BitrateLimitMax {
		totalBitrate = 	settings.Encoding.BitrateLimitMax
	}
	if totalBitrate < settings.Encoding.BitrateLimitMin {
		maxLength := float64(targetSizeKbit) / float64(settings.Encoding.BitrateLimitMin)
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		os.Exit(0)
	}
	// Target select
	metadata.SelectEncoder(totalBitrate)

	// Audio encoding check
	if settings.SelectedAEncoder.UsesBitrate == false {
		log.Println("Encoding audio...")
		audioBitrate = encoder.EncodeAudio(settings.InputVideo) / 1000
	} else {
		audioBitrate = int((float64(settings.SelectedAEncoder.BitratePerc) / float64(100)) * float64(totalBitrate))
		if audioBitrate > settings.SelectedAEncoder.MaxBitrate {
			audioBitrate = settings.SelectedAEncoder.MaxBitrate
		}
		if audioBitrate < settings.SelectedAEncoder.MinBitrate {
			audioBitrate = settings.SelectedAEncoder.MinBitrate
		}
	}

	// Video bitrate calc
	videoBitrate = totalBitrate - audioBitrate
	if settings.Debug {
		log.Println("Calculated target bitrate: " + strconv.Itoa(totalBitrate) + "k")
		log.Println("Audio bitrate: " + strconv.Itoa(audioBitrate) + "k")
		log.Println("Video bitrate: " + strconv.Itoa(videoBitrate) + "k")
	}

	// Encode
	if settings.SelectedAEncoder.UsesBitrate {
		audioMerge = false
	} else {
		audioMerge = true
	}
	if settings.SelectedVEncoder.TwoPass {
		log.Println("Encoding, pass 1/2")
		encoder.Encode(*inputVideo, 1, videoBitrate, audioMerge, audioBitrate)
		log.Println("Encoding, pass 2/2")
		encoder.Encode(*inputVideo, 2, videoBitrate, audioMerge, audioBitrate)
	} else {
		log.Println("Encoding, pass 1/1")
		encoder.Encode(*inputVideo, 0, videoBitrate, audioMerge, audioBitrate)
	}
	log.Println("Cleaning up...")
	os.Remove("ffmpeg2pass-0.log")
	os.Remove("ffmpeg2pass-0.log.mbtree")
	os.Remove(strings.TrimSuffix(settings.InputVideo, path.Ext(settings.InputVideo)) + ".audio." + settings.SelectedVEncoder.Container)
	log.Println("Finished!")
}

func CalculateBitrate(video *settings.VidStats, targetSize int) int{
	Bitrate := float64(targetSize) / video.Duration
	return int(Bitrate)
}