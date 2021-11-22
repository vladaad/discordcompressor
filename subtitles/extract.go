package subtitles

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"log"
	"os"
	"os/exec"
)

func ExtractSubs(filename string, startingTime float64, totalTime float64) string {
	var options []string
	outFilename := utils.GenUUID() + ".mkv"

	options = append(options, "-y", "-loglevel", "error")


	options = append(options, "-i", filename)
	// the times have to be after the filename for subs to work nice
	options = append(options, metadata.AppendTimes(startingTime, totalTime)...)

	options = append(options, "-c", "copy", "-c:s", "ass", "-map", "0", "-vn", "-an")

	options = append(options, outFilename)

	if settings.Debug {
		log.Println(options)
	}

	if !settings.DryRun {
		cmd := exec.Command(settings.General.FFmpegExecutable, options...)

		if settings.ShowStdOut {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			panic(err)
		}
	}

	return outFilename
}
