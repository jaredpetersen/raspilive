package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/jaredpetersen/raspilive/internal/video"
	"github.com/kelseyhightower/envconfig"
)

// TODO create server directory if it does not exist
// TODO research application logging

// HlsSpecification represents the application configuration for the HLS mode.
type HlsSpecification struct {
	Fps *int
}

// Specification represents the application configuration.
type Specification struct {
	Mode string
	Hls  HlsSpecification
}

var errorInvalidSpecification = errors.New("invalid specification")

func main() {
	var specification Specification
	err := envconfig.Process("raspilive", &specification)
	if err != nil {
		log.Fatal(err.Error())
	}

	switch strings.ToUpper(specification.Mode) {
	case "HLS":
		log.Println("Using HLS mode...")
		video.ServeHls()
	default:
		log.Fatal(errorInvalidSpecification)
		os.Exit(1)
	}
}
