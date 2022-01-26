package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
)

func EncodeAudio(video *settings.Video) (outBitrate float64, outFilename string) {
	// filename
	outFilenameBase := video.UUID + "."
	// normalize audio
	lnParams := new(LoudnormParams)
	if video.Input.AudioChannels > 2 {
		video.Output.Audio.Normalize = true
	}
	if video.Output.Audio.Normalize {
		dec := decodeAudio(video, lnParams)
		lnParams = detectVolume(dec)
		if isAudioSilent(lnParams) {
			return 0.0, ""
		}
	}
	// start decoding
	dec := decodeAudio(video, lnParams)
	// encode
	switch video.Output.Audio.Encoder.Type {
	case "ffmpeg":
		outFilename = outFilenameBase + video.Output.Video.Encoder.Container
		encFFmpeg(outFilename, video, dec)
	case "qaac":
		outFilename = outFilenameBase + "m4a"
		encQaac(outFilename, video, dec)
	case "fdkaac":
		outFilename = outFilenameBase + "m4a"
		encFDKaac(outFilename, video, dec)
	default:
		log.Println("Encoder type " + video.Output.Audio.Encoder.Type + " not found")
		os.Exit(0)
	}
	// bitrate
	if !settings.DryRun {
		return metadata.GetStats(outFilename, true).Bitrate, outFilename
	} else {
		return video.Output.Audio.Bitrate, outFilename
	}
}

func isAudioSilent(params *LoudnormParams) bool {
	if params.LRA == "-inf" || params.Thresh == "-inf" || params.IL == "-inf" || params.TP == "-inf" {
		return true
	}
	return false
}
