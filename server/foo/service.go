package foo

import (
	"context"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"github.com/wanmei002/websocket-reverse-proxy/reverse_dialer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	messages.UnimplementedFooServer
	barClient messages.BarClient
}

func New() (*Service, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("passthrough:%s", "device-01"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(reverse_dialer.Dialer()),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		barClient: messages.NewBarClient(conn),
	}, nil
}

func (svc *Service) GetServerAddress(ctx context.Context, in *emptypb.Empty) (*messages.GetAddressResponse, error) {
	return svc.barClient.GetAddress(ctx, in)
}
