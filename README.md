# raspilive
Raspilive is a small application that streams live video from the Raspberry Pi Camera module and converts it into 
various consumable formats.

## Usage
### Configuration
Raspilive relies on environment variables to configure its behavior. You must always specify which video format mode
should be used and then specify any required configurations for those modes (see below). If an optional configuration
value is not set, the default provided by `raspivid` or Ffmpeg will be used instead.

| Environment Variable             | Required | Description                                                     |
| -------------------------------- | -------- | --------------------------------------------------------------- |
| `RASPILIVE_MODE`                 | True     | Streaming video format that is to be used, e.g. HLS, DASH, etc. |
| `RASPILIVE_VIDEO_WIDTH`          | False    | Video width                                                     |
| `RASPILIVE_VIDEO_HEIGHT`         | False    | Video height                                                    |
| `RASPILIVE_VIDEO_FPS`            | False    | Video frames per second                                         |
| `RASPILIVE_VIDEO_HORIZONTALFLIP` | False    | Flip the video horizontally                                     |
| `RASPILIVE_VIDEO_VERTICALFLIP`   | False    | Flip the video vertically                                       |

#### HLS
[HLS](https://en.wikipedia.org/wiki/HTTP_Live_Streaming) is a video streaming format that works by splitting the video
into small consummable segments that are arranged in a continuously changing playlist of files. The client reads from
the playlist and downloads the video segments as needed. DASH requires a static file server to serve all of these files
and raspilive provides this out of the box automatically.

| Environment Variable         | Required | Description                                                            |
| ---------------------------- | -------- | ---------------------------------------------------------------------- |
| `RASPILIVE_HLS_PORT`         | True     | Static file server port number                                         |
| `RASPILIVE_HLS_DIRECTORY`    | False    | Location on disk where files are to be stored and served from          |
| `RASPILIVE_HLS_SEGMENTTYPE`  | False    | Video segment type                                                     |
| `RASPILIVE_HLS_SEGMENTTIME`  | False    | Duration of the video segments in seconds                              |
| `RASPILIVE_HLS_PLAYLISTSIZE` | False    | Maximum number of entries in the playlist at one time                  |
| `RASPILIVE_HLS_STORAGESIZE`  | False    | Maximum number of unreferenced segments to keep on disk before removal |

#### DASH
[DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP), also known as MPEG-DASH, is a video
streaming format that works by splitting the video into small consummable segments that are arranged in a continuously
changing playlist of files. The client reads from the playlist and downloads the video segments as needed. DASH
requires a static file server to serve all of these files and raspilive provides this out of the box automatically.

| Environment Variable         | Required | Description                                                            |
| ---------------------------- | -------- | ---------------------------------------------------------------------- |
| `RASPILIVE_DASH_PORT`         | True     | Static file server port number                                         |
| `RASPILIVE_DASH_DIRECTORY`    | False    | Location on disk where files are to be stored and served from          |
| `RASPILIVE_DASH_SEGMENTTIME`  | False    | Duration of the video segments in seconds                              |
| `RASPILIVE_DASH_PLAYLISTSIZE` | False    | Maximum number of entries in the playlist at one time                  |
| `RASPILIVE_DASH_STORAGESIZE`  | False    | Maximum number of unreferenced segments to keep on disk before removal |

## Installation
Raspilive uses [raspivid](https://www.raspberrypi.org/documentation/usage/camera/raspicam/raspivid.md) to operate the
Raspberry Pi Camera Module. This is already available on the Raspbian operating system and can be enabled via 
[raspi-config](https://www.raspberrypi.org/documentation/configuration/raspi-config.md).

Raspilive also uses [Ffmpeg](https://ffmpeg.org/), a prominent video conversion command line utility, to process the
streaming video that the Raspberry Pi Camera Module outputs. Version 4.0 or higher is required.
```zsh
sudo apt-get install ffmpeg
```
