package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

func main() {
	wsUrl := url.URL{
		Scheme: "ws",
		Host:   "127.0.0.1:8099",
		Path:   "/ws",
	}
	log.Printf("connecting to %s", wsUrl.String())
	c, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		panic(err)
	}
	tk := time.NewTicker(time.Second * 120)
	defer tk.Stop()
	var i int
	for {
		select {
		case <-tk.C:
			log.Printf("time end")
			return
		default:
			err = c.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("hello world %d", i)))
			if err != nil {
				panic(err)
			}
			log.Printf("send %d\n", i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
