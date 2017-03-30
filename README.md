# raspi-live
Stream live video from the Raspberry Pi camera module to the web using HLS

## Dependencies
* FFmpeg
* Raspberry Pi camera module

## Why HLS instead of MPEG DASH?
While MPEG DASH is a more open standard for streaming video over the internet, HLS is a more widely adopted and supported format. FFmpeg, the video conversion software Live Pi uses to process the H.264 video stream from the camera module, technically supports both formats but to varying degrees. The MPEG DASH format is not listed in the [FFmpeg format documentation](https://www.ffmpeg.org/ffmpeg-formats.html) and there are minor performance issues in terms of framerate drops and and degraded image quality. While these performance issues can likely be configured away with more advanced options, HLS is simply easier to implement out of the box.
