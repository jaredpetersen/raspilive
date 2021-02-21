package server

import (
	"log"
	"net/http"
	"strconv"
)

// ServeFiles starts a static file server for everything in the configured directory.
//
// All files are available via the /camera endpoint.
//
// This is a blocking operation.
func ServeFiles(port int, directory string) {
	fs := http.FileServer(http.Dir(directory))
	http.Handle("/camera", fs)

	log.Println("Server started on port", port)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
