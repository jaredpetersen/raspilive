package server

import (
	"errors"
	"net/http"
	"os"
	"strconv"
)

// NewStatic creates a static file server.
//
// This is useful for video streaming options like HLS and DASH which rely on serving video segments. The Directory
// must exist before serving files from it.
func NewStatic(port int, directory string) (*http.Server, error) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return nil, errors.New("directory does not exist")
	}

	router := http.NewServeMux()
	router.Handle("/camera", http.FileServer(http.Dir(directory)))

	server := http.Server{Addr: ":" + strconv.Itoa(port), Handler: router}

	return &server, nil
}
