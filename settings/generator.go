package settings

import (
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func benchmarkx264() float64 {
	log.Println("Testing your PC....")
	log.Println("This may take up to 20 seconds, be patient!")
	cmd := exec.Command(General.FFmpegExecutable,
		"-f", "lavfi", "-i", "nullsrc=192x108", "-vframes", "60",
		"-vf", "geq=random(1)*255:128:128,scale=-2:1080:flags=neighbor",
		"-c:v", "libx264", "-preset", "fast", "-crf", "51", "-f", "null", utils.NullDir(),
		)

	start := time.Now()
	err := cmd.Start()

	err = cmd.Wait()
	elapsed := time.Since(start)

	if err != nil {
		log.Println("Benchmark failed")
		return 0
	}

	score := 1.5 * 60 / elapsed.Seconds()
	log.Println("Benchmark score: " + strconv.FormatFloat(score, 'f', 0, 64))

	return score
}