<!-- DO NOT REMOVE - contributor_list:data:start:["vladaad", "anddddrew"]:end -->

# DiscordCompressor

A small program in Go that efficiently compresses videos using FFmpeg to a certain filesize.

## Dependencies

[FFmpeg](https://ffmpeg.org)

## Usage

`discordcompressor <arguments> <input video(s)>`

* `-o filename` - Sets the output filename, extension is automatically added
* `-size 25` - Sets the target size in MB
* `-last 10` - Compresses the last x seconds of a video
* `-ss 15` - Sets the start time of the video in seconds
* `-t 10` - Sets the time to encode after the start of the file or `-ss` in seconds
* `-mixaudio` - Mixes all audio tracks into one
* `-normaudio` - Normalizes audio volume (use if the input video's audio is very quiet, loud or uneven)
* `-settings string` - Selects the settings file if you have multiple, or generates a new one with the chosen suffix.
* `-debug` - Shows extra information. Please use when reporting bugs, or if you're just curious.
* `-c:v` - Forces a certain video encoder, specified in settings.json
* `-c:a` - Forces a certain audio encoder, specified in settings.json
* `-f` - Forces a certain container, for example, `-f mkv` will output a .mkv file.

Settings and logs are located in %appdata%\vladaad\dc on Windows and ~/.config/vladaad/dc on Linux

## Compiling from source

You need [Go 1.16](https://golang.org/dl/) or newer

Afterwards, run `go build`

<!-- prettier-ignore-start -->
<!-- DO NOT REMOVE - contributor_list:start -->

## 👥 Contributors

- **[@vladaad](https://github.com/vladaad)**

- **[@anddddrew](https://github.com/anddddrew)**

<!-- DO NOT REMOVE - contributor_list:end -->
<!-- prettier-ignore-end -->