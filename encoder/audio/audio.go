package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"path"
	"strings"
)

func EncodeAudio (inFilename string, inBitrate float64, audioTracks int) (outBitrate float64, outFilename string) {
	// filename
	outFilenameBase := strings.TrimSuffix(inFilename, path.Ext(inFilename)) + ".audio."
	// encode
	switch settings.SelectedAEncoder.Type {
	case "ffmpeg":
		outFilename = outFilenameBase + settings.SelectedVEncoder.Container
		encFFmpeg(inFilename, outFilename, inBitrate, audioTracks)
	case "qaac":
		outFilename = encQaac(inFilename, inBitrate, audioTracks)
	default:
		log.Println("Encoder type " + settings.SelectedAEncoder.Type + " not found")
		os.Exit(0)
	}
	// bitrate
	return metadata.GetStats(outFilename, true).Bitrate, outFilename
}
