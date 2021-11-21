package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
)

func EncodeAudio (inFilename string, UUID string, inBitrate float64, container string, eOptions *settings.AudioEncoder, stats *metadata.VidStats, startingTime float64,  totalTime float64) (outBitrate float64, outFilename string) {
	// filename
	outFilenameBase := UUID + "."
	// normalize audio
	lnParams := new(LoudnormParams)
	if settings.Advanced.NormalizeAudio {
		dec := decodeAudio(inFilename, startingTime, totalTime, false, stats, lnParams)
		lnParams = detectVolume(dec)
		if isAudioSilent(lnParams) {
			return 0.0, ""
		}
	}
	// start decoding
	dec := decodeAudio(inFilename, startingTime, totalTime, settings.Advanced.NormalizeAudio, stats, lnParams)
	// encode
	switch eOptions.Type {
	case "ffmpeg":
		outFilename = outFilenameBase + container
		encFFmpeg(outFilename, inBitrate, eOptions, dec)
	case "qaac":
		outFilename = outFilenameBase + "m4a"
		encQaac(outFilename, inBitrate, eOptions, dec)
	case "fdkaac":
		outFilename = outFilenameBase + "m4a"
		encFDKaac(outFilename, inBitrate, eOptions, dec)
	default:
		log.Println("Encoder type " + eOptions.Type + " not found")
		os.Exit(0)
	}
	// bitrate
	if !settings.DryRun {
		return metadata.GetStats(outFilename, true).Bitrate, outFilename
	} else {
		return inBitrate, outFilename
	}
}

func isAudioSilent (params *LoudnormParams) bool {
	if params.LRA == "-inf" || params.Thresh == "-inf" || params.IL == "-inf" || params.TP == "-inf" {
		return true
	}
	return false
}