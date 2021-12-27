package subtitles

import (
	"github.com/vladaad/discordcompressor/metadata"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
)

func ExtractSubs(video *settings.Video) string {
	var options []string
	outFilename := video.UUID + ".subs.mkv"

	options = append(options, "-y", "-loglevel", "error")


	options = append(options, "-i", video.Filename)
	// the times have to be after the filename for subs to work nice
	options = append(options, metadata.AppendTimes(video)...)

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
			return ""
			os.Remove(outFilename)
		}
		err = cmd.Wait()
		if err != nil {
			return ""
			os.Remove(outFilename)
		}
	}

	return outFilename
}
