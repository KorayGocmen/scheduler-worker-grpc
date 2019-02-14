package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	workerID string
)

func init() {
	loadConfig()
}

func main() {

	go startGRPCServer()
	go registerWorker()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-sig:
			fatal(fmt.Sprintf("Signal (%d) received, stopping\n", s))
		}
	}
}

func fatal(message string) {
	deregisterWorker()
	log.Fatalln(message)
}
