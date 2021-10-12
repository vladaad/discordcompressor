package settings

import (
	"encoding/json"
	"fmt"
	"github.com/vladaad/discordcompressor/utils"
	"os"
)

// Stolen from https://github.com/Wieku/danser-go/app/settings

var fileStorage *fileformat
var fileName string

func initStorage() {
	fileStorage = &fileformat{
		General:   General,
		Decoding:  Decoding,
		Encoding:  Encoding,
		Advanced:  Advanced,
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
	defer file.Close()

	if os.IsNotExist(err) {
		saveSettings(fileName, fileStorage)
		return true
	} else if err != nil {
		panic(err)
	} else {
		load(file, fileStorage)
		saveSettings(fileName, fileStorage) //this is done to save additions from the current format
	}
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