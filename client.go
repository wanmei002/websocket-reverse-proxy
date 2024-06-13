package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"os"
	"time"
)

func main() {
	dialer, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		panic(err)
	}

	defer dialer.Close()
	args := os.Args
	if len(args) < 2 {
		panic("please input deviceID")
	}
	first := tcp_server.Device{
		Type: tcp_server.DeviceTypeClient,
		ID:   "",
		To:   args[1],
	}

	sendData, err := json.Marshal(first)
	if err != nil {
		panic(err)
	}
	_, err = dialer.Write(append(sendData, tcp_server.MessageEnd))
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("start conn")
	for i := 0; i < 10; i++ {
		_, err = dialer.Write([]byte(fmt.Sprintf("hello word %d", i)))
		if err != nil {
			panic(err)
		}
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}
