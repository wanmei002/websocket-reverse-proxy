package server

import (
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"github.com/wanmei002/websocket-reverse-proxy/proxy"
	"github.com/wanmei002/websocket-reverse-proxy/server/bar"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var sockPath = "/home/ubuntu/bar.sock"

func BarRun() {
	svr := grpc.NewServer()
	messages.RegisterBarServer(svr, bar.New())
	reflection.Register(svr)
	os.Remove(sockPath)
	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		if err := svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	unixConn, err := net.Dial("unix", sockPath)
	if err != nil {
		log.Fatalf("failed to dial unix: %v", err)
	}

	go func() {
		err = proxy.Run("device-01", unixConn)
		if err != nil {
			log.Fatalf("failed to proxy: %v", err)
		}
	}()

	return

}
