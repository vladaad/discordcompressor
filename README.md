<!-- DO NOT REMOVE - contributor_list:data:start:["vladaad", "notandrewdev"]:end -->
# discordcompressor
A small program in Go that efficiently compresses videos using ffmpeg.

## Dependencies
[FFmpeg](https://ffmpeg.org/) including FFprobe
### Optional (needed for some settings options)
[qaac](https://github.com/nu774/qaac)

## Usage
`discordcompressor <arguments> <input video(s)>`
 * `-focus string` - Sets the focus - for example, "framerate" or "resolution" (configured in settings.json)
 * `-mixaudio` - Mixes all audio tracks into one
 * `-noscale` - Disables FPS limiting and scaling (not recommended)
 * `-reenc string` - Force re-encodes audio or video ("a", "v", "av")
 * `-settings string` - Selects the settings file - for example, settings-test.json.
 * `-size 8` - Sets the target size in MB
 * `-last 10` - Compresses the last x seconds of a video
 * `-ss 15` - Sets the starting time like in ffmpeg
 * `-t 10` - Sets the time to encode after the start of the file or `-ss`

Settings and logs are located in %appdata%\vladaad\dc on Windows and /home/username/.config/vladaad/dc on other platforms.

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