package settings

import (
	"log"
	"os"
	"os/exec"
)

func checkEncoder(encoder string, audio bool) bool {
	var options []string

	if audio {
		options = append(options, "-f", "lavfi", "-i", "anullsrc")
	} else {
		options = append(options, "-f", "lavfi", "-i", "nullsrc=1920x1080")
	}

	options = append(options, "-t", "1")

	if audio {
		options = append(options, "-c:a", encoder)
	} else {
		options = append(options, "-c:v", encoder)
	}

	options = append(options, "-f", "null", "-")

	if Debug {
		log.Println(options)
	}

	// Execution
	cmd := exec.Command(General.FFmpegExecutable, options...)

	if Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println("Testing ", encoder)
	}

	err := cmd.Start()
	if err != nil {
		return false
	}
	err = cmd.Wait()
	if err != nil {
		return false
	}

	return true
}
