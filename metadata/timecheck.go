package metadata

import (
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
)

func CalculateTime(video *settings.Vid, lastSecs float64, targetStart float64, targetTime float64) *settings.Vid {
	video.Time.Duration = video.Input.Duration
	if lastSecs != 0 {
		video.Time.Start = video.Input.Duration - lastSecs
		video.Time.Duration = lastSecs
		return video
	}
	if targetStart != 0 {
		if targetStart > video.Input.Duration {
			log.Println("Starting time is longer than the video!")
			os.Exit(1)
		}
		video.Time.Start = targetStart
		video.Time.Duration = video.Input.Duration - targetStart
	}
	if targetTime != 0 {
		if targetTime > (video.Input.Duration - targetStart) {
			targetTime = video.Input.Duration - targetStart
		}
		video.Time.Duration = targetTime
	}
	return video
}
