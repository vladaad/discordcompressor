package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"strconv"
)

func CalculateTarget(filename string, targetSize float64) *settings.OutTarget {
	target := new(settings.OutTarget)
	// Total bitrate calc
	settings.MaxTotalBitrate = targetSize / settings.VideoStats.Duration
	SelectEncoder(settings.MaxTotalBitrate)
	if settings.MaxTotalBitrate > settings.Encoding.BitrateLimitMax {
		settings.MaxTotalBitrate = 	settings.Encoding.BitrateLimitMax
	}
	if settings.MaxTotalBitrate < settings.Encoding.BitrateLimitMin {
		maxLength := targetSize / settings.Encoding.BitrateLimitMin
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		os.Exit(0)
	}
	// Audio calc
	if settings.VideoStats.AudioTracks != 0 {
		target.AudioBitrate = settings.MaxTotalBitrate * float64(settings.SelectedAEncoder.BitratePerc) / float64(100)
		if target.AudioBitrate > settings.SelectedAEncoder.MaxBitrate {
			target.AudioBitrate = settings.SelectedAEncoder.MaxBitrate
		}
		if target.AudioBitrate < settings.SelectedAEncoder.MinBitrate {
			target.AudioBitrate = settings.SelectedAEncoder.MinBitrate
		}
	}
	target.AudioPassthrough, target.VideoPassthrough, target.AudioBitrate = CheckStreamCompatibility(filename, target.AudioBitrate)
	if target.AudioPassthrough {
		target.AudioBitrate = settings.VideoStats.AudioBitrate
	}
	// Encode audio
	if !target.AudioPassthrough && settings.VideoStats.AudioTracks != 0 {
		log.Println("Encoding audio...")
		target.AudioBitrate = EncodeAudio(filename, target.AudioBitrate)
	}
	if settings.VideoStats.AudioTracks != 0 {
		if settings.Debug {
			log.Println("Audio passthrough: "  + strconv.FormatBool(target.AudioPassthrough))
			if !target.AudioPassthrough {
				log.Println("Target audio bitrate: " + strconv.FormatFloat(target.AudioBitrate, 'f', 1, 64) + "k")
			}
		}
		target.VideoBitrate = settings.MaxTotalBitrate - target.AudioBitrate
	} else {
		target.VideoBitrate = settings.MaxTotalBitrate
		target.AudioBitrate = -1
	}
	if settings.Debug {
		log.Println("Video passthrough: "  + strconv.FormatBool(target.VideoPassthrough))
		if !target.VideoPassthrough {
			log.Println("Target video bitrate: " + strconv.FormatFloat(target.VideoBitrate, 'f', 1, 64) + "k")
		}
	}
	return target
}
