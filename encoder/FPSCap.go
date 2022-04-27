package encoder

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func calculateFPS(video *settings.Vid) *settings.Vid {
	video.Output.FPS = video.Input.FPS
	for {
		if float64(video.Output.FPS.N)/float64(video.Output.FPS.D) > float64(video.Output.Settings.MaxFPS) {
			if settings.Encoding.HalveFPS {
				video.Output.FPS.D *= 2
			} else {
				video.Output.FPS.N = video.Output.Settings.MaxFPS
				video.Output.FPS.D = 1
				break
			}
		} else {
			break
		}
	}
	return video
}

func fpsFilter(video *settings.Vid) []string {
	fps := video.Output.FPS
	if video.Input.FPS != video.Output.FPS {
		var str []string
		expr := strconv.Itoa(fps.N)
		expr += "/"
		expr += strconv.Itoa(fps.D)
		str = append(str, "-r", expr)
		return str
	} else {
		return nil
	}
}
