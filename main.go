package main

import (
	"github.com/wanmei002/websocket-reverse-proxy/server"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//go func() {
	//	unixs.UnixListener()
	//}()
	go func() {
		err := tcp_server.Run()
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)
	//unixDial, err := proxy.DialUnix()
	//if err != nil {
	//	panic(err)
	//}
	//err = proxy.Run("device-01", unixDial)
	//if err != nil {
	//	panic(err)
	//}

	server.FooRun()
	server.BarRun()

	signals := make(chan os.Signal, 3)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:
		log.Println("Got shutdown signal")
	}
}
