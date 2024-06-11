package main

import (
	"github.com/wanmei002/websocket-reverse-proxy/websockets"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", websockets.Websocket)
	err := http.ListenAndServe(":8099", nil)
	if err != nil {
		panic(err)
	}
}
