package audio

import (
	"bufio"
	"encoding/json"
	"github.com/vladaad/discordcompressor/settings"
	"github.com/vladaad/discordcompressor/utils"
	"io"
	"os/exec"
	"strings"
)

type LoudnormParams struct {
	IL     string `json:"input_i"`
	LRA    string `json:"input_lra"`
	TP     string` json:"input_tp"`
	Thresh string `json:"input_thresh"`
}

func detectVolume(input io.ReadCloser) *LoudnormParams {
	var options []string
	// input
	options = append(options, "-i", "-")
	// filter
	options = append(options, "-af", "loudnorm=print_format=json")
	// out
	options = append(options, "-f", "null", utils.NullDir())

	cmd := exec.Command(settings.General.FFmpegExecutable, options...)

	cmd.Stdin = input

	pipe, err := cmd.StderrPipe()
	scanner := bufio.NewScanner(pipe)
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	// prepare eye protection
	foundJson := false
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line,"{") {
			foundJson = true
		}
		if foundJson {
			lines = append(lines, line)
		}
		if strings.Contains(line,"}") {
			foundJson = false
		}
	}
	rawJson := strings.Join(lines, "\n")

	parsed := new(LoudnormParams)
	err = json.Unmarshal([]byte(rawJson), &parsed)
	if err != nil {
		panic(err)
	}

	return parsed
}