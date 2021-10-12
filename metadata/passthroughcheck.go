package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
)

func CheckStreamCompatibility(filename string, audioBitrateIn float64, bitrate float64, videoStats *VidStats, startingTime float64, totalTime float64, vEncoder *settings.Encoder, aEncoder *settings.AudioEncoder) (audioCompatible bool, videoCompatible bool, audioBitrateOut float64) {
	audioCompatible, videoCompatible = false, false
	// If bitrate wasn't able to be analyzed, analyze it
	if (videoStats.AudioBitrate == 0 || videoStats.VideoBitrate == 0) && videoStats.AudioTracks != 0 {
		videoStats.AudioBitrate = AnalyzeAudio(filename)
	}
	format := findCurrentFormat(vEncoder.Container)
	audioFmt := findACodec(videoStats.AudioCodec, format)
	videoFmt := findVCodec(videoStats.VideoCodec, format)
	// VB calc
	if videoStats.AudioTracks != 0 {
		videoStats.VideoBitrate = videoStats.Bitrate - videoStats.AudioBitrate
	} else {
		videoStats.VideoBitrate = videoStats.Bitrate
	}
	// To save you from understanding this spaghetti, the audio has to be:
	// The same codec as would normally be encoded
	// Below max bitrate
	if audioFmt != nil && videoStats.AudioBitrate < aEncoder.MaxBitrate && videoStats.AudioTracks != 0 {
		audioCompatible = true
		audioBitrateIn = videoStats.AudioBitrate
	}

	// The conditions for video compatibility:
	// The same codec as would normally be encoded
	// Video bitrate must be detected
	// Audio should be passed through too
	// Video bitrate must be below total bitrate
	if videoFmt != nil && utils.Contains(videoStats.Pixfmt, videoFmt.PixFmts) {
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

// yes there are better ways to do this but I am lazy

func findCurrentFormat(container string) *settings.Format {
	for i := range settings.Advanced.CompatibleFormats {
		if settings.Advanced.CompatibleFormats[i].Container == container {
			return settings.Advanced.CompatibleFormats[i]
		}
	}
	return nil
}

func findVCodec(codec string, format *settings.Format) *settings.VideoCodec {
	for i := range format.CompatibleVideoCodecs {
		if format.CompatibleVideoCodecs[i].Name == codec {
			return format.CompatibleVideoCodecs[i]
		}
	}
	return nil
}

func findACodec(codec string, format *settings.Format) *settings.AudioCodec {
	for i := range format.CompatibleAudioCodecs {
		if format.CompatibleAudioCodecs[i].Name == codec {
			return format.CompatibleAudioCodecs[i]
		}
	}
	return nil
}