package settings

import (
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func populateSettings() {
	Encoding.AudioEncoders = []*AudioEncoder{generateAudioEncoder()}
	selectPresets()
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
	} else if !strings.Contains(utils.CommandOutput("ffmpeg", "-h encoder=libfdk_aac"), "is not recognized") {
		encoder = &AudioEncoder{
			Name:         "aac",
			Type:         "ffmpeg",
			Encoder:      "libfdk_aac",
			CodecName:    "aac",
			Options:      "-cutoff 17500",
			UsesBitrate:  true,
			MaxBitrate:   128,
			MinBitrate:   96,
			BitratePerc:  10,
		}
	} else if utils.CheckIfPresent("fdkaac") {
		encoder = &AudioEncoder{
			Name:         "aac",
			Type:         "fdkaac",
			Encoder:      "",
			CodecName:    "aac",
			Options:      "-w 17500",
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
			MaxBitrate:   192,
			MinBitrate:   160,
			BitratePerc:  10,
		}
		// use twoloop if possible
		if strings.Contains(utils.CommandOutput("ffmpeg", "-h encoder=aac"), "twoloop") {
			encoder.Options = "-aac_coder twoloop"
			encoder.MaxBitrate = 160
			encoder.MinBitrate = 112
		}
	}
	return encoder
}

func selectPresets() {
	presets := []string{"veryfast", "faster", "fast", "medium", "slow", "slower", "veryslow"}
	offsets := []int{2, 3, 4, 4, 5, 6, 6} // look into encoding.go - offset from the "base" fast preset
	offset := 0
	slowest := 6
	fastest := 0
	score := 0.0

	if ForceScore == -1 {
		score = benchmarkx264()
	} else {
		score = ForceScore
	}

	// Select offsets depending on score
	// yes, this is a mess
	if score > 120 { // score 120+ - medium is fastest, medium -> slow,...
		offset = 1
		fastest = 4
	} else if score > 70 { // score 70+ - fast is fastest, medium -> slow,...
		offset = 1
	} else if score > 45 { // score 45+ - default
		offset = 0
	} else if score > 30 { // score 30+ - slower is slowest, medium -> fast,...
		offset = -1
		slowest = 5
	} else if score > 20 { // score 20+ - faster is fastest, slow is slowest, medium -> faster,...
		offset = -2
		fastest = 1
	} else if score > 10 { // score 10+ - medium -> veryfast,...
		offset = -3
	} else { // potato - medium is slowest, medium -> veryfast,...
		offset = -3
		slowest = 3
	}

	// Apply presets
	for i := range offsets {
		presetN := offsets[i] + offset

		// holy mother of spaghetti code
		presetN = int(math.Min(float64(presetN), float64(slowest))) // clamp to slowest
		presetN = int(math.Max(float64(presetN), float64(fastest))) // clamp to fastest

 		Encoding.BitrateTargets[i].Preset = presets[presetN]
	}
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
		return 50
	}

	score := 1.5 * 60 / elapsed.Seconds()
	log.Println("Benchmark score: " + strconv.FormatFloat(score, 'f', 0, 64))

	return score
}