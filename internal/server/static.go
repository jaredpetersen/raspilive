package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Static is a static file server.
//
// Files may be accessed via the route `/camera`.
type Static struct {
	Port      int    // Port the server runs on. Uses the next available port if one is not provided.
	Cert      string // Location of a certificate file for TLS
	Key       string // Location of a key file for TLS
	Directory string // Directory the files should be served from
	listener  net.Listener
	server    http.Server
}

// ListenAndServe begins listening on the configured port and serving static files.
//
// ListenAndServe always returns a non-nil error. When the server closes, this will be http.ErrServerClosed.
func (stcsrv *Static) ListenAndServe() error {
	if err := stcsrv.listen(); err != nil {
		return err
	}

	return stcsrv.serve()
}

func (stcsrv *Static) listen() error {
	// If a Port is not chose, listener will choose the next available port for us
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(stcsrv.Port))
	if err != nil {
		return err
	}

	stcsrv.listener = listener
	stcsrv.Port = listener.Addr().(*net.TCPAddr).Port

	return nil
}

func (stcsrv *Static) serve() error {
	_, err := os.Stat(stcsrv.Directory)
	if os.IsNotExist(err) {
		return errors.New("directory does not exist")
	}

	router := http.NewServeMux()
	router.Handle("/camera/", http.StripPrefix("/camera", http.FileServer(http.Dir(stcsrv.Directory))))

	log.Println("Server started on port", stcsrv.Port)
	if stcsrv.Cert != "" && stcsrv.Key != "" {
		err = stcsrv.server.ServeTLS(stcsrv.listener, stcsrv.Cert, stcsrv.Key)
	} else {
		stcsrv.server = http.Server{Handler: router}
		err = stcsrv.server.Serve(stcsrv.listener)
	}

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the server with a deadline.
//
// Gives active connections the opportunity to finish their work within the given time period before ultimately
// closing all connections and shutting everything down.
func (stcsrv *Static) Shutdown(timout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timout)
	defer cancel()

	return stcsrv.server.Shutdown(ctx)
}
