package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/vladaad/discordcompressor/utils"
)

// Stolen from https://github.com/Wieku/danser-go/app/settings

var fileStorage *fileformat

func initStorage() {
	fileStorage = &fileformat{
		General:  General,
		Decoding: Decoding,
		Encoding: Encoding,
		Advanced: Advanced,
	}
}

func LoadSettings(version string) bool {
	initStorage()
	fileName := utils.SettingsDir()
	fileName += "/settings"
	if version != "" {
		fileName += "-" + version
	}
	fileName += ".json"

	file, err := os.Open(fileName)

	if errors.Is(err, fs.ErrNotExist) {
		populateSettings()
		saveSettings(fileName, fileStorage)
		return true
	} else if err != nil {
		panic(err)
	} else {
		load(file, fileStorage)
		saveSettings(fileName, fileStorage) //this is done to save additions from the current format
	}

	defer file.Close()
	return false
}

func load(file *os.File, target interface{}) {
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(target); err != nil {
		panic(fmt.Sprintf("Failed to parse %s! Please re-check the file for mistakes. Error: %s", file.Name(), err))
	}
}

func saveSettings(path string, source interface{}) {
	file, err := os.Create(path)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(source); err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}
