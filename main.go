package main

import (
	"github.com/wanmei002/websocket-reverse-proxy/proxy"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"github.com/wanmei002/websocket-reverse-proxy/unixs"
	"time"
)

func main() {
	go func() {
		unixs.UnixListener()
	}()
	go func() {
		err := tcp_server.Run()
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(5 * time.Second)
	err := proxy.Run()
	if err != nil {
		panic(err)
	}
}
