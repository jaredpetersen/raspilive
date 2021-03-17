package main

import (
	"os"
	"os/signal"
)

func osStopper(stop chan struct{}) {
	// Set up a channel for OS signals so that we can quit gracefully if the user terminates the program
	// Once we get this signal, sent a message to the stop channel
	osStop := make(chan os.Signal, 1)
	signal.Notify(osStop, os.Interrupt, os.Kill)

	go func() {
		<-osStop
		stop <- struct{}{}
	}()
}
