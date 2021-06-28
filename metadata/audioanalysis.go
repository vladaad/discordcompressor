package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
)

func AnalyzeAudio(filename string) int {
	log.Println("Extracting audio for analysis...")
	extract := exec.Command(
		settings.General.FFmpegExecutable,
		"-i", filename,
		"-map", "0:a:0",
		"-c", "copy",
		"analyze" + filename,
	)

	err := extract.Start()
	if err != nil {
		panic("Couldn't start FFmpeg")
	}

	err = extract.Wait()
	if err != nil {
		panic(err)
	}

	bitrate := GetStats("analyze" + filename).Bitrate
	os.Remove("analyze" + filename)
	return bitrate
}