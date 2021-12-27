package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
)

func CheckStreamCompatibility(video *settings.Video) *settings.Video {
	// If bitrate wasn't able to be analyzed, analyze it
	if (video.Input.AudioBitrate == 0 || video.Input.VideoBitrate == 0) && video.Input.AudioTracks != 0 {
		video.Input.AudioBitrate = AnalyzeAudio(video.Filename)
	}
	format := findCurrentFormat(video.Output.Video.Encoder.Container)
	audioFmt := findACodec(video.Input.AudioCodec, format)
	videoFmt := findVCodec(video.Input.VideoCodec, format)
	// VB calc
	if video.Input.VideoBitrate != 0 {
		video.Input.VideoBitrate = video.Input.Bitrate - video.Input.AudioBitrate
	} else {
		video.Input.VideoBitrate = video.Input.Bitrate
	}
	// To save you from understanding this spaghetti, the audio has to be:
	// The same codec as would normally be encoded
	// Below max bitrate
	if audioFmt != nil && video.Input.AudioBitrate < video.Output.Audio.Encoder.MaxBitrate && video.Input.AudioTracks != 0 {
		video.Output.Audio.Passthrough = true
		video.Output.Audio.Bitrate = video.Input.AudioBitrate
	}

	// The conditions for video compatibility:
	// The same codec as would normally be encoded
	// Video bitrate must be detected
	// Audio should be passed through too
	// Video bitrate must be below total bitrate
	if videoFmt != nil && utils.Contains(video.Input.Pixfmt, videoFmt.PixFmts) {
		if video.Output.Audio.Passthrough && video.Output.TotalBitrate > video.Input.Bitrate {
			video.Output.Video.Passthrough = true
		} else if video.Input.VideoBitrate < video.Output.TotalBitrate - video.Output.Audio.Bitrate {
			video.Output.Video.Passthrough = true
		}
	}

	// I'm not dealing with times and passthrough, fuck that
	if video.Time.Time != video.Input.Duration || video.Time.Start != float64(0) {
		video.Output.Video.Passthrough, video.Output.Audio.Passthrough = false, false
	}

	if video.Output.Audio.Mix || video.Output.Audio.Normalize {
		video.Output.Audio.Passthrough = false
	}

	return video
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