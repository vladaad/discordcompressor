package encoder

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func calculateFPS(video *settings.Vid) (*settings.FPS, bool) {
	modified := false // This variable is for accurate detection if FPS was changed, float fuckery? was making -vf fps=25/1 when input was 25.0000
	fps := new(settings.FPS)
	fps.D, fps.N = video.Input.FPS.D, video.Input.FPS.N
	for {
		if float64(fps.N)/float64(fps.D) > float64(video.Output.Settings.MaxFPS)+1 { // allow for leniency
			if settings.Encoding.HalveFPS {
				fps.D *= 2
				modified = true
			} else {
				fps.N = video.Output.Settings.MaxFPS
				fps.D = 1
				modified = true
				break
			}
		} else {
			break
		}
	}
	return fps, modified
}

func fpsFilter(video *settings.Vid) string {
	fps := video.Output.FPS
	if video.Input.FPS != video.Output.FPS {
		var str string
		expr := strconv.Itoa(fps.N)
		expr += "/"
		expr += strconv.Itoa(fps.D)
		str = "fps=" + expr
		return str
	} else {
		return ""
	}
}
