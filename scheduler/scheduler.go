package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	loadConfig()
}

func main() {

	go api()
	go startGRPCServer()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}
}
