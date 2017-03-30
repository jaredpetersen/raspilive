# raspi-live
Stream live video from the Raspberry Pi camera module to the web using HLS

## Dependencies
Raspi-live uses FFmpeg, a video conversion command-line utlity, to process the streaming H.264 video that the Raspberry Pi camera module outputs. The [Hannes ihm sein Blog](http://hannes.enjoys.it/blog/2016/03/ffmpeg-on-raspbian-raspberry-pi/) has a great guide on how to install FFmpeg on the Raspberry Pi if you are not familiar with the process.

The raspivid command-line tool is also required, though it should be installed automatically as part of the [normal camera module installation process](https://www.raspberrypi.org/documentation/usage/camera/).

## Video Stream Playback
Raspi-live only handles streaming live video from the camera module and is not concerned with how the stream is played. The streaming content is in HLS format and can be passed to any HLS-capable player via the `/camera/camera1.m3u8` resource.

## Why HLS instead of MPEG DASH?
While MPEG DASH is a more open standard for streaming video over the internet, HLS is more widely adopted and supported. FFmpeg, the video conversion software Live Pi uses to process the H.264 video stream from the camera module, technically supports both formats but to varying degrees. The MPEG DASH format is not listed in the [FFmpeg format documentation](https://www.ffmpeg.org/ffmpeg-formats.html) and there are minor performance issues in terms of framerate drops and and degraded image quality. While these performance issues can likely be configured away with more advanced options, HLS is simply easier to implement out of the box.
