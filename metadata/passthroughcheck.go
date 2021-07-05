package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
)

func CheckStreamCompatibility(filename string, audioBitrateIn float64) (audioCompatible bool, videoCompatible bool, audioBitrateOut float64) {
	audioCompatible, videoCompatible = false, false
	// If bitrate wasn't able to be analyzed, analyze it xd
	if (settings.VideoStats.AudioBitrate == 0 || settings.VideoStats.VideoBitrate == 0) && settings.VideoStats.AudioTracks != 0 {
		settings.VideoStats.AudioBitrate = AnalyzeAudio(filename)
	}
	// VB calc
	if settings.VideoStats.AudioTracks != 0 {
		settings.VideoStats.VideoBitrate = settings.VideoStats.Bitrate - settings.VideoStats.AudioBitrate
	} else {
		settings.VideoStats.VideoBitrate = settings.VideoStats.Bitrate
	}
	// To save you from understanding this spaghetti, the audio has to be:
	// The same codec as would normally be encoded
	// Below max bitrate
	if settings.VideoStats.AudioCodec == settings.SelectedAEncoder.CodecName && settings.VideoStats.AudioBitrate < settings.SelectedAEncoder.MaxBitrate && settings.VideoStats.AudioTracks != 0 {
		audioCompatible = true
		audioBitrateIn = settings.VideoStats.AudioBitrate
	}

	if settings.VideoStats.AudioTracks > 1 {
		log.Println("Multiple audio tracks detected! You can use -mixaudio to mix them into one")
	}

	// The conditions for video compatibility:
	// The same codec as would normally be encoded
	// Video bitrate must be detected
	// Audio should be passed through too
	// Video bitrate must be below total bitrate
	if settings.SelectedVEncoder.CodecName == settings.VideoStats.VideoCodec && (settings.VideoStats.Pixfmt == "yuv420p" || settings.VideoStats.Pixfmt == "yuv420p10le" || settings.VideoStats.Pixfmt == settings.SelectedVEncoder.Pixfmt) {
		if audioCompatible && settings.MaxTotalBitrate < settings.VideoStats.Bitrate {
			videoCompatible = true
		} else if settings.VideoStats.VideoBitrate < settings.MaxTotalBitrate - audioBitrateIn {
			videoCompatible = true
		}
	}

	// I'm not dealing with times and passthrough, fuck that
	if settings.Time != float64(0) || settings.Starttime != float64(0) {
		audioCompatible, videoCompatible = false, false
	}

	return audioCompatible, videoCompatible, audioBitrateIn
}