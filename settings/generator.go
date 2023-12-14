package settings

import (
	"github.com/vladaad/discordcompressor/utils"
	"log"
)

func newSettings() {
	log.Println("Generating new settings")
	genAudioEncoders()
	genVideoEncoders()
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

func genVideoEncoders() {
	fastest := genFastestEncoder()

	Encoding.Encoders = []*Encoder{
		fastest,
		{
			Name:   "fast",
			Passes: 2,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v libx264 -preset medium -aq-mode 3",
		},
		{
			Name:   "normal",
			Passes: 2,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v libx264 -preset slow -aq-mode 3",
		},
		{
			Name:   "slow",
			Passes: 2,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v libx264 -preset veryslow -aq-mode 3",
		},
		{
			Name:   "veryslow",
			Passes: 2,
			Keyint: 15,
			Pixfmt: "yuv420p10le",
			Args:   "-c:v libvpx-vp9 -lag-in-frames 25 -cpu-used 4 -auto-alt-ref 1 -arnr-maxframes 7 -arnr-strength 4 -aq-mode 0 -enable-tpl 1 -row-mt 1", // credit: BlueSwordM
		},
		{
			Name:   "ultra",
			Passes: 2,
			Keyint: 15,
			Pixfmt: "yuv420p10le",
			Args:   "-c:v libaom-av1 -cpu-used 6 -lag-in-frames 25 -aom-params enable-keyframe-filtering=1:arnr-strength=1:sb-size=dynamic:enable-chroma-deltaq=1:enable-qm=1:quant-b-adapt=1", // credit: wiki.x266.mov/blog/av1-encoding-for-dummies
		},
	}
}

func genFastestEncoder() *Encoder {
	var encoder *Encoder
	log.Println("Checking hardware encoders, this may take a while")
	// fuck this shit i hate multithreading
	nvenc := checkEncoder("h264_nvenc", false)
	qsv := checkEncoder("h264_qsv", false)
	amf := checkEncoder("h264_amf", false)

	if nvenc {
		encoder = &Encoder{
			Name:   "fastest",
			Passes: 1,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v h264_nvenc -preset p7",
		}
	} else if qsv {
		encoder = &Encoder{
			Name:   "fastest",
			Passes: 1,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v h264_qsv -preset veryslow",
		}
	} else if amf {
		encoder = &Encoder{
			Name:   "fastest",
			Passes: 1,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v h264_amf -quality quality",
		}
	} else {
		encoder = &Encoder{
			Name:   "fastest",
			Passes: 1,
			Keyint: 10,
			Pixfmt: "yuv420p",
			Args:   "-c:v libx264 -preset veryfast -aq-mode 3",
		}
	}
	return encoder
}

func genAACEncoder() *AudioEncoder {
	var encoder *AudioEncoder

	if utils.CheckIfPresent("qaac64") {
		encoder = &AudioEncoder{
			Name:  "aac",
			Type:  "qaac",
			BMult: 1.2,
			BMax:  144,
			BMin:  72,
			TVBR:  false,
			Args:  "",
		}
	} else if utils.CheckIfPresent("fhgaacenc") {
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
