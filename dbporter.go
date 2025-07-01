package main

import (
	"Mehul-Kumar-27/dbporter/logger"
	"os"
	"os/signal"
	"syscall"
)

type DbPorter struct {
	logger logger.Logger
}

func StartPorter() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	done := make(chan struct{})

	go func() {

	}()

	select {
	case <-done:

	}
}
