package settings

import (
	"github.com/vladaad/discordcompressor/utils"
	"log"
)

func newSettings() {
	log.Println("Generating new settings")
	genAudioEncoders()
}

func genAudioEncoders() {
	aac := genAACEncoder()
	opus := genOpusEncoder()

	if Debug {
		log.Println("Automatic encoder choice")
		log.Println("AAC: ", aac.Type, ", Opus: ", opus.Type)
	}

	Encoding.AEncoders = []*AudioEncoder{aac, opus}
}

func genAACEncoder() *AudioEncoder {
	var encoder *AudioEncoder
	if utils.CheckIfPresent("fhgaacenc") {
		encoder = &AudioEncoder{
			Name:  "aac",
			Type:  "fhgaac",
			BMult: 1.2,
			BMax:  144,
			BMin:  72,
			TVBR:  false,
			Args:  "",
		}
	} else if utils.CheckIfPresent("fdkaac") {
		encoder = &AudioEncoder{
			Name:  "aac",
			Type:  "fdkaac",
			BMult: 1.2,
			BMax:  144,
			BMin:  72,
			TVBR:  false,
			Args:  "",
		}
	} else {
		encoder = &AudioEncoder{
			Name:  "aac",
			Type:  "ffmpeg",
			BMult: 1.4,
			BMax:  192,
			BMin:  96,
			TVBR:  false,
			Args:  "-c:a aac",
		}
	}
	return encoder
}

func genOpusEncoder() *AudioEncoder {
	var encoder *AudioEncoder // for future expansion
	encoder = &AudioEncoder{
		Name:  "opus",
		Type:  "ffmpeg",
		BMult: 1,
		BMax:  128,
		BMin:  32,
		TVBR:  false,
		Args:  "-c:a libopus",
	}
	return encoder
}
