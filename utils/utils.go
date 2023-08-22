package utils

import (
	"github.com/vladaad/discordcompressor/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

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

func CheckIfPresent(filename string) bool {
	_, err := exec.Command(filename).Output()
	if err == nil {
		return true
	}
	return !strings.Contains(err.Error(), "executable file not found")
}

func GetArg(args string, argToFind string) string {
	split1 := strings.Split(args, argToFind)
	if len(split1) < 2 {
		return ""
	}
	split2 := strings.Split(split1[1], " ")
	return split2[1]
}
