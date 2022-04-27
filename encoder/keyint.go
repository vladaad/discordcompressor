package encoder

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
)

func keyint(video *settings.Vid) string {
	return strconv.Itoa(video.Output.Encoder.Keyint * video.Output.FPS.N / video.Output.FPS.D)
}
