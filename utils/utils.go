package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func OpenURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func GenUUID() string {
	raw := uuid.New()
	cleaned := strings.ReplaceAll(raw.String(), "-", "")
	return cleaned
}

func Contains(input string, list []string) bool {
	for i := range list {
		if input == list[i] {
			return true
		}
	}
	return false
}

func ContainsInt(input int, list []int) bool {
	for i := range list {
		if input == list[i] {
			return true
		}
	}
	return false
}