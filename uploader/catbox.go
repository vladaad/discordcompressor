package uploader

import (
	"github.com/wabarc/go-catbox"
	"log"
	"net/http"
	"time"
)

func catboxUpload(path string) string {

	client := &http.Client{
		Timeout: 1200 * time.Second,
	}

	url, err := catbox.New(client).Upload(path)

	if err != nil {
		log.Println(err)
		panic("Catbox upload failed!")
	}

	return url
}
