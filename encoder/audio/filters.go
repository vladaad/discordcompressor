package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"strconv"
	"strings"
)

func filters(mixTracks bool, normalize bool, videoStats *metadata.VidStats, lnParams *LoudnormParams) (filter string, mapping string) {
	var filters []string
	var inputs []string
	if mixTracks {
		for i := 0; i < videoStats.AudioTracks; i++ {
			inputs = append(inputs, "[0:a:" + strconv.Itoa(i) + "]")
		}
	} else {
		inputs = []string{"[0:a:0]"}
	}
	mapping = "0:a:0"
	if mixTracks {
		var filter []string
		filter = append(filter, inputs...)
		filter = append(filter, "amix=inputs=", strconv.Itoa(videoStats.AudioTracks))
		filter = append(filter, "[mixed]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[mixed]"}
		mapping = inputs[0]
	}
	if normalize {
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