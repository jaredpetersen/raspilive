# raspi-live
raspi-live is a Node.js Express webserver that takes streaming video from the Raspberry Pi Camera module and makes it available on the web via [HLS](https://en.wikipedia.org/wiki/HTTP_Live_Streaming) or [DASH](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP).

Run it via a simple command line interface:
```
raspi-live start
```

The server will start serving the streaming files on `/camera`. Point streaming video players to `/camera/livestream.m3u8` for HLS or `/camera/livestream.mpd` for DASH.


## Usage
```
$ raspi-live --help

  Usage: raspi-live [options] [command]

  self-contained raspberry pi video streaming server

  Options:

    -v, --version    output the version number
    -h, --help       output usage information

  Commands:

    start [options]  start streaming video from the raspberry pi camera module
```

### Options
#### -v, --version
Output the version number.

#### -h, --help
Output information on how to use the command line interface.

### Commands
#### start \[options\]
Start streaming video from the raspberry pi camera module.

##### Options
###### -d, --directory
The directory used to host the streaming video files. Those concerned about the long-term health of their pi's SD card may opt to point raspi-live to a RAMDisk so that the files are only stored in memory. However, this also means that you will be unable to recover any of the footage if the power is cut.

Defaults to `/home/<USERNAME>/camera` but `/srv/camera` is recommended as raspi-live is a server.

###### -f, --format
* [`hls`](https://en.wikipedia.org/wiki/HTTP_Live_Streaming) (default)
* [`dash`](https://en.wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP)

###### -l, --list-size
The number of streaming video files included in the playlist.

Defaults to `10`.

###### -s, --storage-size
The number of streaming video files stored after they cycle out of the playlist. This is useful in cases where you want to look at previously recorded footage. The streaming video files are 2 seconds long so to have a 24-hour cycle of recorded video, specify `43200` (make sure to have enough storage space).

Defaults to `10`.

###### -p, --port
Port number the server runs on.

Defaults to `8080`.


## Install
### Raspberry Pi Camera module
raspi-live only supports streaming video from the Raspberry Pi camera module. Here's a the official documentation on how to connect and configure it: https://www.raspberrypi.org/documentation/usage/camera/.

### FFmpeg
raspi-live uses FFmpeg, a video conversion command-line utility, to process the streaming H.264 video that the Raspberry Pi camera module outputs. Here's how to install it on your Raspberry Pi:

1. Download and configure FFmpeg via:
```
sudo apt-get install libomxil-bellagio-dev
wget -O ffmpeg.tar.bz2 https://ffmpeg.org/releases/ffmpeg-snapshot-git.tar.bz2
tar xvjf ffmpeg.tar.bz2
cd ffmpeg
sudo ./configure --arch=armel --target-os=linux --enable-gpl --enable-omx --enable-omx-rpi --enable-nonfree
```
2. If you're working with a Raspbery Pi 2 or 3, then run `sudo make -j4` to build FFmpeg. If you're working with a Raspberry Pi Zero, then run `sudo make`.
3. Install FFmpeg via `sudo make install` regardless of the model of your Raspberry Pi.
4. Delete the FFmpeg directory and tar file that were created during the download process in Step 1. FFmpeg has been installed so they are no longer needed.

### CLI
Install it globally:
```
npm install raspi-live -g
raspi-live --help
```
Or use npx:
```
npx raspi-live --help
```


## Video Stream Playback
raspi-live is only concerned with streaming video from the camera module and does not offer a playback solution.

Browser support between the different streaming formats varies so in most cases a JavaScript playback library will be necessary. For more information on this, check out [Mozilla's article on the subject](https://developer.mozilla.org/en-US/docs/Web/Apps/Fundamentals/Audio_and_video_delivery/Live_streaming_web_audio_and_video).
