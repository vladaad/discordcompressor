package audio

import (
	"github.com/vladaad/discordcompressor/metadata"
	"strconv"
	"strings"
)

func filters(mixTracks bool, videoStats *metadata.VidStats) (filter string, mapping string) {
	var filters []string
	var inputs []string
	for i := 0; i < videoStats.AudioTracks; i++ {
		inputs = append(inputs, "[0:a:" + strconv.Itoa(i) + "]")
	}
	mapping = "0:a:0"
	if mixTracks {
		var filter []string
		filter = append(filter, inputs...)
		filter = append(filter, "amix=inputs=", strconv.Itoa(videoStats.AudioTracks))
		filter = append(filter, "[mixed]")
		filters = append(filters, strings.Join(filter, ""))
		inputs = []string{"[mixed]"}
	}

	merged := strings.Join(filters, ";")
	return merged, mapping
}