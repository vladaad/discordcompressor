package uploader

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
)

func Upload(path string) string {
	// Checking if filesize is under limit
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if file.Size() > int64(settings.General.UploadMaxMB*1048576) { // 1024^2
		log.Println("The video has exceeded the set upload limit! Please reduce the size parameter.")
		return ""
	}

	// Upload
	log.Println("Uploading...")
	switch settings.General.UploadService {
	case "catbox":
		return catboxUpload(path)
	default:
		log.Println("Upload service name invalid!")
		return ""
	}
}