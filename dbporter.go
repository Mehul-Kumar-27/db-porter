package main

import (
	"os"
	"os/signal"
	"syscall"
)

type DbPorter struct {
}

func StartPorter() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {

	}()

	select {
	case <-done:

	}
}
