package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"strconv"
)

func CalculateTarget(filename string, targetSize int) *settings.OutTarget {
	target := new(settings.OutTarget)
	// Total bitrate calc
	settings.MaxTotalBitrate = int(float64(targetSize) / settings.VideoStats.Duration)
	SelectEncoder(settings.MaxTotalBitrate)
	if settings.MaxTotalBitrate > settings.Encoding.BitrateLimitMax {
		settings.MaxTotalBitrate = 	settings.Encoding.BitrateLimitMax
	}
	if settings.MaxTotalBitrate < settings.Encoding.BitrateLimitMin {
		maxLength := float64(targetSize) / float64(settings.Encoding.BitrateLimitMin)
		log.Println("File too long! Maximum length: " + strconv.FormatFloat(maxLength, 'f', 1, 64) + " seconds")
		os.Exit(0)
	}
	// Audio calc
	if settings.VideoStats.AudioTracks != 0 {
		target.AudioBitrate = int(float64(settings.MaxTotalBitrate) * (float64(settings.SelectedAEncoder.BitratePerc) / 100))
		if target.AudioBitrate > settings.SelectedAEncoder.MaxBitrate {
			target.AudioBitrate = settings.SelectedAEncoder.MaxBitrate
		}
		if target.AudioBitrate < settings.SelectedAEncoder.MinBitrate {
			target.AudioBitrate = settings.SelectedAEncoder.MinBitrate
		}
	}
	target.AudioPassthrough, target.VideoPassthrough = CheckStreamCompatibility(filename)
	if target.AudioPassthrough {
		target.AudioBitrate = settings.VideoStats.AudioBitrate
	}
	// Encode audio
	if !target.AudioPassthrough {
		log.Println("Encoding audio...")
		target.AudioBitrate = EncodeAudio(filename)
	}
	if settings.VideoStats.AudioTracks != 0 {
		if settings.Debug {
			log.Println("Audio passthrough: "  + strconv.FormatBool(target.AudioPassthrough))
			if !target.AudioPassthrough {
				log.Println("Target audio bitrate: " + strconv.Itoa(target.AudioBitrate) + "k")
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
			log.Println("Target video bitrate: " + strconv.Itoa(target.VideoBitrate) + "k")
		}
	}
	return target
}
