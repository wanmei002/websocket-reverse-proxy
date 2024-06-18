package server

import (
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"github.com/wanmei002/websocket-reverse-proxy/server/bar"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func BarRun() {
	UnaryServerInterceptorOtelp("bar")
	svr := grpc.NewServer(
		//grpc.UnaryInterceptor(UnaryServerInterceptorOtelp("bar")),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	messages.RegisterBarServer(svr, bar.New())
	reflection.Register(svr)
	lis, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		fmt.Println("bar run")
		if err := svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return

}
