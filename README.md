# raspilive
Raspilive is a small command-line application that streams live video from the Raspberry Pi Camera module to the web

## Usage
```
raspilive is a livestreaming tool for the Raspberry Pi Camera Module

For more information visit https://github.com/jaredpetersen/raspilive

Usage:
  raspilive [command]

Available Commands:
  hls         Stream video using HLS
  dash        Stream video using DASH
  help        Help about any command

Flags:
  -h, --help      help for raspilive
  -v, --version   version for raspilive
```

### HLS
[HLS](https://en.wikipedia.org/wiki/HTTP_Live_Streaming) is a video streaming format that works by splitting the video
into small consummable segments that are arranged in a continuously changing playlist of files. The client reads from
the playlist and downloads the video segments as needed. DASH requires a static file server to serve all of these files
and raspilive provides this out of the box automatically.

#### Flags
##### --width int
Video width. Defaults to 1920.

##### --height int
Video height. Defaults to 1080.

##### --fps int
Video framerate. Defaults to 30.

##### --horizontal-flip
Horizontally flip the video.

##### --vertical-flip
Vertically flip the video.

##### --port int
Static file server port. Finds an available port if one is not specified.

##### --directory string
Static file server directory. Defaults to the current directory.

##### --tls-cert string
Static file server TLS certificate.

##### --tls-key string
Static file server TLS key.

##### --segment-time int
Target segment duration in seconds.

##### --playlist-size int
Maximum number of playlist entries.

##### --storage-size int
Maximum number of unreferenced segments to keep on disk before removal.

### DASH
[DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP), also known as MPEG-DASH, is a video
streaming format that works by splitting the video into small consummable segments that are arranged in a continuously
changing playlist of files. The client reads from the playlist and downloads the video segments as needed. DASH
requires a static file server to serve all of these files and raspilive provides this out of the box automatically.

#### Flags
##### --width int
Video width. Defaults to 1920.

##### --height int
Video height. Defaults to 1080.

##### --fps int
Video framerate. Defaults to 30.

##### --horizontal-flip
Horizontally flip the video.

##### --vertical-flip
Vertically flip the video.

##### --port int
Static file server port. Finds an available port if one is not specified.

##### --directory string
Static file server directory. Defaults to the current directory.

##### --tls-cert string
Static file server TLS certificate.

##### --tls-key string
Static file server TLS key.

##### --segment-time int
Target segment duration in seconds.

##### --playlist-size int
Maximum number of playlist entries.

##### --storage-size int
Maximum number of unreferenced segments to keep on disk before removal.

## Installation
Raspilive uses [raspivid](https://www.raspberrypi.org/documentation/usage/camera/raspicam/raspivid.md) to operate the
Raspberry Pi Camera Module. This is already available on the Raspbian operating system and can be enabled via 
[raspi-config](https://www.raspberrypi.org/documentation/configuration/raspi-config.md).

Raspilive also uses [Ffmpeg](https://ffmpeg.org/), a prominent video conversion command line utility, to process the
streaming video that the Raspberry Pi Camera Module outputs. Version 4.0 or higher is required.
```zsh
sudo apt-get install ffmpeg
```

Download the latest version of Raspilive from the [Releases page](https://github.com/jaredpetersen/raspi-live/releases).
All of the release binaries are compiled for ARM 6 and are compatible with Raspberry Pi.
