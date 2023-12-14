package metadata

import (
	"encoding/json"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type StreamList struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

type Stream struct {
	CodecName  string `json:"codec_name"`
	StreamType string `json:"codec_type"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Pixfmt     string `json:"pix_fmt"`
	Framerate  string `json:"r_frame_rate"`
	Channels   int    `json:"channels"`
	Bitrate    string `json:"bit_rate"`
}

type Format struct {
	Duration string `json:"duration"`
	Bitrate  string `json:"bit_rate"`
}

func GetStats(filename string, audioonly bool) *settings.InputStats {
	stats := new(settings.InputStats)
	stats.Bitrate = new(settings.Bitrates)
	probe, err := exec.Command(
		settings.General.FFprobeExecutable,
		"-loglevel", "quiet",
		"-of", "json",
		"-show_entries", "stream:format",
		filename,
	).Output()

	if err != nil {
		panic("Failed to start FFprobe")
	}
	// JSON parsing
	probeOutput := new(StreamList)

	err = json.Unmarshal(probe, &probeOutput)
	if err != nil {
		log.Println(err)
		panic("Failed to parse JSON")
	}

	if !audioonly {
		videoStream := findFirstStream("video", probeOutput.Streams)
		stats.Bitrate.Video, _ = strconv.Atoi(videoStream.Bitrate)
		stats.VCodec = videoStream.CodecName
		stats.Pixfmt = videoStream.Pixfmt
		stats.Width = videoStream.Width
		stats.Height = videoStream.Height
		// FPS
		fps := new(settings.FPS)
		fpsSplit := strings.Split(videoStream.Framerate, "/")
		n1, _ := strconv.Atoi(fpsSplit[0])
		n2, _ := strconv.Atoi(fpsSplit[1])
		fps.N, fps.D = n1, n2
		stats.FPS = fps
	}
	// Other
	stats.Duration, _ = strconv.ParseFloat(probeOutput.Format.Duration, 64)
	stats.Bitrate.Total, _ = strconv.Atoi(probeOutput.Format.Bitrate)
	// Audio
	stats.ATracks = countStreams("audio", probeOutput.Streams)
	if stats.ATracks != 0 {
		audioStream := findFirstStream("audio", probeOutput.Streams)
		stats.ACodec = audioStream.CodecName
		stats.Bitrate.Audio, _ = strconv.Atoi(audioStream.Bitrate)
		stats.AChannels = audioStream.Channels
	}

	return stats
}

func findFirstStream(streamType string, streamList []Stream) Stream {
	for i := range streamList {
		if streamList[i].StreamType == streamType {
			return streamList[i]
		}
	}
	panic("Stream not found")
}

func countStreams(streamType string, streamList []Stream) int {
	streamCount := 0
	for i := range streamList {
		if streamList[i].StreamType == streamType {
			streamCount++
		}
	}
	return streamCount
}
