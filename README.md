<!-- DO NOT REMOVE - contributor_list:data:start:["vladaad", "notandrewdev"]:end -->
# discordcompressor
A small program in Go that efficiently compresses videos using ffmpeg.

## Dependencies
[FFmpeg](https://ffmpeg.org/) including FFprobe

## Usage
`discordcompressor <arguments>`
 * `-debug` - Prints extra info
 * `-focus string` - Sets the focus - for example, "framerate" or "resolution" (configured in settings.json)
 * `-i path` - Sets the input video
 * `-mixaudio` - Mixes all audio tracks into one
 * `-noscale` - Disables FPS limiting and scaling
 * `-settings string` - Selects the settings file - for example, settings-test.json
 * `-size 8` - Sets the target size in MB
 * `-ss 15` - Sets the starting time like in ffmpeg
 * `-t 10` - Sets the time to encode after the start of the file or `-ss`
 
## Compiling from source
You need [Go 1.16](https://golang.org/dl/) or newer

Afterwards run `go build` or `build.bat`

<!-- prettier-ignore-start -->
<!-- DO NOT REMOVE - contributor_list:start -->
## 👥 Contributors


- **[@vladaad](https://github.com/vladaad)**

- **[@notandrewdev](https://github.com/notandrewdev)**

<!-- DO NOT REMOVE - contributor_list:end -->
<!-- prettier-ignore-end -->