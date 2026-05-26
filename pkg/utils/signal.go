package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func HandleSignals(function func()) {
	var callback sync.Once
	sigc := make(chan os.Signal, 1)
	defer close(sigc)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigc
	callback.Do(function)
}
