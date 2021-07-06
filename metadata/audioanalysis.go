package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strings"
)

func AnalyzeAudio(filename string) float64 {
	log.Println("Extracting audio for analysis...")
	sFilename := strings.Split(filename, ".")
	extension := sFilename[len(sFilename)-1]
	outFilename := strings.Replace(filename, extension, "analyse." + extension, len(sFilename) - 1)
	extract := exec.Command(
		settings.General.FFmpegExecutable,
		"-i", filename,
		"-map", "0:a:0",
		"-c", "copy",
		outFilename,
	)
	err := extract.Start()
	if err != nil {
		panic("Couldn't start FFmpeg")
	}

	err = extract.Wait()
	if err != nil {
		panic(err)
	}

	bitrate := GetStats(outFilename, true).Bitrate
	os.Remove(outFilename)
	return bitrate
}