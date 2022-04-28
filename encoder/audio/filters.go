package audio

import (
	"github.com/vladaad/discordcompressor/settings"
	"strconv"
	"strings"
)

// dynaudnorm=g=5:f=300:p=0.99

func getFilters(video *settings.Vid) (filter string, mapping string) {
	var filters []string
	var inputs []string
	// Input
	if settings.MixAudio {
		for i := 0; i < video.Input.ATracks; i++ {
			inputs = append(inputs, "[0:a:"+strconv.Itoa(i)+"]")
		}
	} else {
		inputs = append(inputs, "[0:a:0]")
	}
	mapping = "0:a:0"
	// Audio mixing
	if settings.MixAudio {
		var filter []string
		filter = append(filter, inputs...)
		filter = append(filter, "amix=inputs=", strconv.Itoa(video.Input.ATracks))
		filter = append(filter, "[mixed]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[mixed]"}
		mapping = inputs[0]
	}

	filter = strings.Join(filters, ";")
	return filter, mapping
}
