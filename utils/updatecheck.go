package utils

import (
	"encoding/json"
	"github.com/vladaad/discordcompressor/build"
	"log"
	"net/http"
	"strings"
	"time"
)

// Stolen from danser-go

func CheckForUpdates() {
	if strings.Contains(build.VERSION, "dev") { //false positive, those are changed during compile
		return
	}

	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/vladaad/discordcompressor/releases/latest", nil)
	if err != nil {
		log.Println("Failed to check for updates - Can't create request")
		return
	}

	client := new(http.Client)
	response, err := client.Do(request)

	if err != nil || response.StatusCode != 200 {
		log.Println("Failed to check for updates - Can't get release info from GitHub")
		return
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	var data struct {
		URL string `json:"html_url"`
		Tag string `json:"tag_name"`
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Println("Failed to check for updates - Failed to decode the response from GitHub")
		return
	}

	if data.Tag != build.VERSION {
		log.Println("You're using an older version of discordcompressor.")
		log.Println("You can download a newer version here:", data.URL)
		time.Sleep(time.Second * 2)
	} else {
		log.Println("Discordcompressor is up to date!")
	}
}