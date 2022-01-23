package metadata

import (
	"encoding/json"
	"github.com/vladaad/discordcompressor/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type StreamList struct {
	Frames  []Frame  `json:"frames"`
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
	Samplerate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
	Bitrate    string `json:"bit_rate"`
	Tags       Tag    `json:"tags"`
}

type Format struct {
	Duration string `json:"duration"`
	Bitrate  string `json:"bit_rate"`
}

type Frame struct {
	SideDataList []SideData `json:"side_data_list"`
}

type SideData struct {
	Type string `json:"side_data_type"`
}

type Tag struct {
	Language string `json:"language""`
}

func GetStats(filepath string, audioonly bool) *settings.VidStats {
	stats := new(settings.VidStats)
	if _, err := os.Stat(filepath); err != nil {
		panic(filepath + " doesn't exist")
	}

	probe, err := exec.Command(
		settings.General.FFprobeExecutable,
		"-loglevel", "quiet",
		"-of", "json",
		"-show_entries", "stream:format",
		"-show_frames", "-read_intervals", "%+#1",
		filepath,
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
		stats.VideoBitrate, _ = strconv.ParseFloat(videoStream.Bitrate, 64)
		stats.VideoCodec = videoStream.CodecName
		stats.Pixfmt = videoStream.Pixfmt
		stats.Width = videoStream.Width
		stats.Height = videoStream.Height
		// FPS
		fpsSplit := strings.Split(videoStream.Framerate, "/")
		n1, _ := strconv.ParseFloat(fpsSplit[0], 64)
		n2, _ := strconv.ParseFloat(fpsSplit[1], 64)
		stats.FPS = n1 / n2
		// HDR detect
		if len(probeOutput.Frames) != 0 {
			if len(probeOutput.Frames[0].SideDataList) != 0 {
				for i := range probeOutput.Frames[0].SideDataList {
					if probeOutput.Frames[0].SideDataList[i].Type == "Mastering display metadata" {
						stats.IsHDR = true
					}
				}
			}
		}
	}
	// Other
	stats.Duration, _ = strconv.ParseFloat(probeOutput.Format.Duration, 64)
	stats.Bitrate, _ = strconv.ParseFloat(probeOutput.Format.Bitrate, 64)
	// Audio
	stats.AudioTracks = countStreams("audio", probeOutput.Streams)
	if stats.AudioTracks != 0 {
		audioStream := findFirstStream("audio", probeOutput.Streams)
		stats.AudioCodec = audioStream.CodecName
		stats.AudioBitrate, _ = strconv.ParseFloat(audioStream.Bitrate, 64)
		stats.SampleRate, _ = strconv.Atoi(audioStream.Samplerate)
		stats.AudioChannels = audioStream.Channels
	}
	// Subtitles
	for i := range probeOutput.Streams {
		if probeOutput.Streams[i].StreamType == "subtitle" {
			if probeOutput.Streams[i].Tags.Language == settings.Advanced.SubfinderLang && !stats.MatchingSubs {
				stats.MatchingSubs = true
				stats.SubtitleStream = i
			}
		}
	}

	// Bitrates -> k
	stats.Bitrate /= 1024
	stats.VideoBitrate /= 1024
	stats.AudioBitrate /= 1024
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