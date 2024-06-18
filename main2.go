package main

import (
	"github.com/wanmei002/websocket-reverse-proxy/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server.BarRun()

	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:
		log.Println("Got shutdown signal")
	}
}
