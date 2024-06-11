package websockets

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{}

func Websocket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("websocket upgrade failed;err:", err.Error())
		return
	}
	go handleWebsocketConn(conn)
}

func handleWebsocketConn(conn *websocket.Conn) {
	wg := sync.WaitGroup{}
	defer conn.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			mt, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("read message failed;err:", err.Error())
				return
			}
			fmt.Println("message type:", mt)
			fmt.Println("message data:", string(data))
		}

	}()

	wg.Wait()
}
