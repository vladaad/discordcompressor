package uploader

import (
	"github.com/wabarc/go-catbox"
	"log"
)

func catboxUpload(path string) string {
	url, err := catbox.New(nil).Upload(path)

	if err != nil {
		log.Println(err)
		panic("Catbox upload failed!")
	}

	return url
}
