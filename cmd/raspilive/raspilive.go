package main

import (
	"github.com/jaredpetersen/raspilive/internal/video"
)

// TODO create server directory if it does not exist
// TODO error handling for raspivid and ffmpeg
// TODO research application logging
// TODO research config (env / YAML) (https://github.com/kelseyhightower/envconfig)

func main() {
	video.ServeHls()
}
