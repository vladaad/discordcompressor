package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

func filters(video *settings.Video, lnParams *LoudnormParams) (filter string, mapping string) {
	var filters []string
	var inputs []string
	if video.Output.Audio.Mix {
		for i := 0; i < video.Input.AudioTracks; i++ {
			inputs = append(inputs, "[0:a:" + strconv.Itoa(i) + "]")
		}
	} else {
		inputs = []string{"[0:a:0]"}
	}
	mapping = "0:a:0"
	if video.Output.Audio.Mix {
		var filter []string
		filter = append(filter, inputs...)
		filter = append(filter, "amix=inputs=", strconv.Itoa(video.Input.AudioTracks))
		filter = append(filter, "[mixed]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[mixed]"}
		mapping = inputs[0]
	} else if video.Input.AudioChannels > 2 { // this is intentional, otherwise downmix is done "normally" via -ac 2
		var filter []string
		filter = append(filter, inputs...)
		switch video.Input.AudioChannels {
		case 6:
			filter = append(filter, "pan=stereo|FL=0.5*FC+0.707*FL+0.707*BL+0.5*LFE|FR=0.5*FC+0.707*FR+0.707*BR+0.5*LFE")
		case 8:
			filter = append(filter, "pan=stereo|FL= 0.5*FC+0.3*FLC+0.3*FL+0.3*BL+0.3*SL+0.5*LFE|FR=0.5*FC+0.3*FRC+0.3*FR+0.3*BR+0.3*SR+0.5*LFE")
		default:
			filter = nil
		}
		filter = append(filter, "[downmixed]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[downmixed]"}
		mapping = inputs[0]
	}
	if video.Output.Audio.Normalize && lnParams.IL != "" {
		var filter []string
		filter = append(filter, inputs...)
		filter = append(filter, "loudnorm=linear=true:i=-14:lra=7:tp=-1")

		filter = append(filter, ":measured_i=" + lnParams.IL)
		filter = append(filter, ":measured_lra=" + lnParams.LRA)
		filter = append(filter, ":measured_tp=" + lnParams.TP)
		filter = append(filter, ":measured_thresh=" + lnParams.Thresh)

		filter = append(filter, "[voladj]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[voladj]"}
		mapping = inputs[0]
	}

	merged := strings.Join(filters, ";")
	return merged, mapping
}