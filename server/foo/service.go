package foo

import (
	"context"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	messages.UnimplementedFooServer
	barClient messages.BarClient
}

func New() (*Service, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("passthrough:%s", "127.0.0.1:9003"),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		barClient: messages.NewBarClient(conn),
	}, nil
}

func (svc *Service) GetServerAddress(ctx context.Context, in *emptypb.Empty) (*messages.GetAddressResponse, error) {
	fmt.Printf("foo ctx: %+v\n", ctx)
	fmt.Printf("foo ctx traceparent: %+v\n", ctx.Value("traceparent"))
	metadata.FromIncomingContext()
	fmt.Println("foo trace:", trace.SpanFromContext(ctx).SpanContext().TraceID().String())
	return svc.barClient.GetAddress(ctx, in)
}
