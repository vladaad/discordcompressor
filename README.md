<!-- DO NOT REMOVE - contributor_list:data:start:["vladaad", "anddddrew"]:end -->
# DiscordCompressor
A small program in Go that efficiently compresses videos using ffmpeg.

## Dependencies
[FFmpeg](https://ffmpeg.org) including FFprobe
### Optional (needed for some settings options)
[qaac](https://github.com/nu774/qaac)
[fdkaac](https://github.com/nu774/fdkaac)

## Usage
`discordcompressor <arguments> <input video(s)>`
 * `-o filename` - Sets the output filename, extension is automatically added
 * `-focus string` - Sets the focus - for example, "framerate" or "resolution" (configured in settings.json)
 * `-mixaudio` - Mixes all audio tracks into one
 * `-normalize` - Normalizes audio volume
 * `-noscale` - Disables FPS limiting and scaling (not recommended)
 * `-reenc string` - Force re-encodes audio or video ("a", "v", "av")
 * `-settings string` - Selects the settings file - for example, settings-test.json.
 * `-forcescore 60` - Forces a specific benchmark score when generating settings. Higher = slower, but higher-quality settings.
 * `-size 8` - Sets the target size in MB
 * `-subfind string` - Finds a certain string in subtitles and cuts according to it
 * `-last 10` - Compresses the last x seconds of a video
 * `-ss 15` - Sets the starting time like in ffmpeg
 * `-t 10` - Sets the time to encode after the start of the file or `-ss`

Settings and logs are located in %appdata%\vladaad\dc on Windows and /home/username/.config/vladaad/dc on other platforms.

## Compiling from source
You need [Go 1.16](https://golang.org/dl/) or newer

Afterwards run `go build` or `build.sh`. `build.sh` builds execs for both 64bit and 32bit and both Windows and Linux.

<!-- prettier-ignore-start -->
<!-- DO NOT REMOVE - contributor_list:start -->
## 👥 Contributors


- **[@vladaad](https://github.com/vladaad)**

- **[@anddddrew](https://github.com/anddddrew)**

<!-- DO NOT REMOVE - contributor_list:end -->
<!-- prettier-ignore-end -->