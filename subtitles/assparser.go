package subtitles

import (
	"strconv"
	"strings"
)

type Line struct {
	startTime float64
	endTime float64
	text string
}

func parseASS(input string) []*Line {
	var textPos int
	var startPos int
	var endPos int
	var startLine int
	split := strings.Split(input, "\n")
	// Get format
	eventsFound := false
	for i := range split {
		line := split[i]
		if strings.Contains(line, "[Events]") {
			eventsFound = true
			startLine = i+2
		}
		if strings.HasPrefix(line, "Format:") && eventsFound {
			clearedPrefix := strings.ReplaceAll(line, "Format: ", "")
			cleanedSpaces := strings.ReplaceAll(clearedPrefix, " ", "")
			format := strings.Split(cleanedSpaces, ",")
			for i := range format {
				if strings.Contains(format[i], "Text"){
					textPos = i
				}
				if strings.Contains(format[i], "Start") {
					startPos = i
				}
				if strings.Contains(format[i], "End") {
					endPos = i
				}
			}
		}
	}
	if !eventsFound {
		return nil
	}
	// Start parsing text
	var lines []*Line
	for i := range split {
		if i >= startLine {
			var out []string
			line := new(Line)
			formatting := false
			splitLine := strings.Split(split[i], ",")
			if len(splitLine) > textPos {
				//parsing characters
				var temp []string
				for i := range splitLine {
					if i >= textPos {
						temp = append(temp, splitLine[i])
					}
				}
				line.text = strings.Join(temp, ",")
				// start & end times
				line.startTime = formatTime(splitLine[startPos])
				line.endTime = formatTime(splitLine[endPos])

			}
			cleaned := strings.ReplaceAll(line.text, "\\N", "")
			chars := strings.Split(cleaned, "")
			for i := range chars {
				// Ignore formatting
				if chars[i] == "{" {
					formatting = true
				}
				if !formatting {
					out = append(out, chars[i])
				}
				if chars[i] == "}" {
					formatting = false
				}
			}
			line.text = strings.Join(out, "")
			lines = append(lines, line)
		}
	}

	return lines
}

func formatTime(time string) (timeInSeconds float64) {
	split := strings.Split(time, ":")
	hours, _ := strconv.ParseFloat(split[0], 64)
	minutes, _ := strconv.ParseFloat(split[1], 64)
	seconds, _ := strconv.ParseFloat(split[2], 64)

	timeInSeconds += hours * 3600
	timeInSeconds += minutes * 60
	timeInSeconds += seconds

	return timeInSeconds
}