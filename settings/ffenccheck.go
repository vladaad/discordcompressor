package settings

import (
	"log"
	"os"
	"os/exec"
)

func CheckAEncoder(encoder string) bool {
	var options []string

	options = append(options, "-f", "lavfi", "-i", "anullsrc")
	options = append(options, "-t", "10")
	options = append(options, "-c:a", encoder)
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
