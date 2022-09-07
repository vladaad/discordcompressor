package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/vladaad/discordcompressor/build"
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
		log.Println("You're lucky that the rickroll didn't work")
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

func NullDir() string {
	var null string
	switch runtime.GOOS {
	case "windows":
		null = "NUL"
	default:
		null = "/dev/null"
	}
	return null
}

func SettingsDir() string {
	if build.BUILD == "portable" {
		fpath, _ := os.Executable()
		return filepath.Dir(fpath)
	} else {
		switch runtime.GOOS {
		case "windows":
			return os.Getenv("APPDATA") + "\\vladaad\\dc"
		default:
			home, _ := os.UserHomeDir()
			return home + "/.config/vladaad/dc"
		}
	}
}

func CommandOutput(cmd string, args []string) string {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
