package settings

import (
	"github.com/vladaad/discordcompressor/utils"
	"strings"
)

func populateSettings() {
	Encoding.AudioEncoders = []*AudioEncoder{generateAudioEncoder()}
}

func generateAudioEncoder() *AudioEncoder {
	var encoder *AudioEncoder
	if utils.CheckIfPresent("qaac64") {
		encoder = &AudioEncoder{
			Name:         "aac",
			Type:         "qaac",
			Encoder:      "",
			CodecName:    "aac",
			Options:      "",
			UsesBitrate:  true,
			MaxBitrate:   128,
			MinBitrate:   96,
			BitratePerc:  10,
		}
	} else {
		encoder = &AudioEncoder{
			Name:         "aac",
			Type:         "ffmpeg",
			Encoder:      "aac",
			CodecName:    "aac",
			Options:      "",
			UsesBitrate:  true,
			MaxBitrate:   160,
			MinBitrate:   128,
			BitratePerc:  10,
		}
		// use twoloop if possible
		if strings.Contains(utils.CommandOutput("ffmpeg", "-h encoder=aac"), "twoloop") {
			encoder.Options = "-aac_coder twoloop"
		}
	}
	return encoder
}