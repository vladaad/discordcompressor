package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"sync"
)

func GenFilename(video *settings.Vid) string {
	filename := video.UUID + "."
	switch video.Output.AEncoder.Type {
	case "ffmpeg":
		filename += video.Output.Settings.Container
	case "qaac", "fdkaac":
		filename += "m4a"
	}
	return filename
}
func EncodeAudio(video *settings.Vid, wg *sync.WaitGroup) *settings.Vid {
	defer wg.Done()
	dec := decodeAudio(video)
	switch video.Output.AEncoder.Type {
	case "ffmpeg":
		encFFmpeg(video, dec)
	case "fdkaac":
		encFDK(video, dec)
	}

	if !settings.Encoding.FastMode {
		video.Output.Bitrate.Audio = getBitrate(video)
		video.Output.Bitrate.Video = video.Output.Bitrate.Total - video.Output.Bitrate.Audio
	} else {
		log.Println("Audio encoding finished!")
	}
	return video
}

func getBitrate(video *settings.Vid) int {
	info, err := os.Stat(video.Output.AudioFile)
	if err != nil {
		panic("Failed to get audio bitrate")
	}
	return int(float64(info.Size()*8) / video.Time.Duration)
}
