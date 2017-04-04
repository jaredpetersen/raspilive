# raspi-live
Stream live video from the Raspberry Pi camera module to the web using HLS

## Starting the Live Stream
Install all of the raspi-live dependencies and then run the usual `npm start`. The application will start streaming encrypted video from the Raspberry Pi camera module via the `/camera/livestream.m3u8` resource on port 8080.

## Dependencies
Raspi-live uses FFmpeg, a video conversion command-line utility, to process the streaming H.264 video that the Raspberry Pi camera module outputs. The [Hannes ihm sein Blog](http://hannes.enjoys.it/blog/2016/03/ffmpeg-on-raspbian-raspberry-pi/) has a great guide on how to install FFmpeg on the Raspberry Pi if you are not familiar with the process.

The raspivid command-line tool is also required, though it should be installed automatically as part of the [normal camera module installation process](https://www.raspberrypi.org/documentation/usage/camera/).

There are also the usual Node.js dependencies that can be installed via `npm install`.

## Video Stream Playback
Raspi-live is only concerned with streaming video from the camera module and is not opinionated about how the stream should be played. Any HLS-capable player should be compatible.

However, the [raspi-live-dash](https://github.com/jaredpetersen/raspi-live-dash) project is available for those looking for an opinion.

## Why HLS instead of MPEG DASH?
While MPEG DASH is a more open standard for streaming video over the internet, HLS is more widely adopted and supported. FFmpeg, chosen for its near-ubiquitousness in the video processing space, technically supports both formats but to varying degrees. The MPEG DASH format is not listed in the [FFmpeg format documentation](https://www.ffmpeg.org/ffmpeg-formats.html) and there are minor performance issues in terms of framerate drops and and degraded image quality with FFmpeg's implementation. While these performance issues can likely be configured away with more advanced options, HLS is simply easier to implement out of the box.
