package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
)

func EncodeAudio (inFilename string, UUID string, inBitrate float64, audioTracks int, container string, eOptions *settings.AudioEncoder, startingTime float64, totalTime float64) (outBitrate float64, outFilename string) {
	// filename
	outFilenameBase := UUID + "."
	// encode
	switch eOptions.Type {
	case "ffmpeg":
		outFilename = outFilenameBase + container
		encFFmpeg(inFilename, outFilename, inBitrate, audioTracks, eOptions, startingTime, totalTime)
	case "qaac":
		outFilename = encQaac(inFilename, outFilenameBase, inBitrate, audioTracks, eOptions, startingTime, totalTime)
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
