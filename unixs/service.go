package unixs

import (
	"fmt"
	"net"
	"os"
)

var SocketFile = "/home/ubuntu/websocket-reverse-proxy.sock"

func UnixListener() {
	os.Remove(SocketFile)
	socketListener, err := net.Listen("unix", SocketFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on "+SocketFile, " successfully")
	for {
		conn, err := socketListener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("New unix connection from " + conn.RemoteAddr().String())
		go func() {
			defer conn.Close()
			for {
				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					panic(err)
				}
				fmt.Println("unix: ", string(buf))
			}
		}()
	}
}
