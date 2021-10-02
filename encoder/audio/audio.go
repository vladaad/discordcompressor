package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"path"
	"strings"
)

func EncodeAudio (inFilename string, inBitrate float64, audioTracks int, startingTime float64, totalTime float64) (outBitrate float64, outFilename string) {
	// filename
	outFilenameBase := strings.TrimSuffix(inFilename, path.Ext(inFilename)) + ".audio."
	// encode
	switch settings.SelectedAEncoder.Type {
	case "ffmpeg":
		outFilename = outFilenameBase + settings.SelectedVEncoder.Container
		encFFmpeg(inFilename, outFilename, inBitrate, audioTracks, startingTime, totalTime)
	case "qaac":
		outFilename = encQaac(inFilename, inBitrate, audioTracks, startingTime, totalTime)
	default:
		log.Println("Encoder type " + settings.SelectedAEncoder.Type + " not found")
		os.Exit(0)
	}
	// bitrate
	return metadata.GetStats(outFilename, true).Bitrate, outFilename
}
