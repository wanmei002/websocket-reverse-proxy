package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"github.com/wanmei002/websocket-reverse-proxy/server/foo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

func FooRun() {
	svr := grpc.NewServer()
	fooSvc, err := foo.New()
	if err != nil {
		panic(err)
	}
	messages.RegisterFooServer(svr, fooSvc)
	reflection.Register(svr)

	lis, err := net.Listen("tcp", ":21113")
	if err != nil {
		log.Fatalf("foo failed to listen: %v", err)
	}

	grpcConn, err := grpc.NewClient("127.0.0.1:21113", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	mux := runtime.NewServeMux()
	err = messages.RegisterFooHandlerClient(context.Background(), mux, messages.NewFooClient(grpcConn))
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		if err := svr.Serve(lis); err != nil {
			log.Fatalf("foo failed to serve: %v", err)
		}
	}()

	httpSvc := &http.Server{
		Addr:    ":21114",
		Handler: mux,
	}
	go func() {
		err = httpSvc.ListenAndServe()
		if err != nil {
			log.Fatalf("foo failed to listen: %v", err)
		}
	}()

	return
}
