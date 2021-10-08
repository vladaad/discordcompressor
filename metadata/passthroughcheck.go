package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
)

func CheckStreamCompatibility(filename string, audioBitrateIn float64, bitrate float64, videoStats *VidStats, startingTime float64, totalTime float64, vEncoder *settings.Encoder, aEncoder *settings.AudioEncoder) (audioCompatible bool, videoCompatible bool, audioBitrateOut float64) {
	audioCompatible, videoCompatible = false, false
	// If bitrate wasn't able to be analyzed, analyze it xd
	if (videoStats.AudioBitrate == 0 || videoStats.VideoBitrate == 0) && videoStats.AudioTracks != 0 {
		videoStats.AudioBitrate = AnalyzeAudio(filename)
	}
	// VB calc
	if videoStats.AudioTracks != 0 {
		videoStats.VideoBitrate = videoStats.Bitrate - videoStats.AudioBitrate
	} else {
		videoStats.VideoBitrate = videoStats.Bitrate
	}
	// To save you from understanding this spaghetti, the audio has to be:
	// The same codec as would normally be encoded
	// Below max bitrate
	if videoStats.AudioCodec == aEncoder.CodecName && videoStats.AudioBitrate < aEncoder.MaxBitrate && videoStats.AudioTracks != 0 {
		audioCompatible = true
		audioBitrateIn = videoStats.AudioBitrate
	}

	// The conditions for video compatibility:
	// The same codec as would normally be encoded
	// Video bitrate must be detected
	// Audio should be passed through too
	// Video bitrate must be below total bitrate
	if vEncoder.CodecName == videoStats.VideoCodec && (videoStats.Pixfmt == "yuv420p" || videoStats.Pixfmt == vEncoder.Pixfmt) {
		if audioCompatible && bitrate < videoStats.Bitrate {
			videoCompatible = true
		} else if videoStats.VideoBitrate < bitrate - audioBitrateIn {
			videoCompatible = true
		}
	}

	// I'm not dealing with times and passthrough, fuck that
	if totalTime != float64(0) || startingTime != float64(0) {
		audioCompatible, videoCompatible = false, false
	}

	return audioCompatible, videoCompatible, audioBitrateIn
}