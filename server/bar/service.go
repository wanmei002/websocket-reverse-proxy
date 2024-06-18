package bar

import (
	"context"
	"errors"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/gen/golang/wanmei002/messages/v1"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type Service struct {
	messages.UnimplementedBarServer
}

func New() *Service {
	return &Service{}
}

func (svc *Service) GetAddress(ctx context.Context, in *emptypb.Empty) (*messages.GetAddressResponse, error) {
	fmt.Println("bar trace:", trace.SpanFromContext(ctx).SpanContext().TraceID().String())
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	ip := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}

	if ip != "" {
		return &messages.GetAddressResponse{Address: ip}, nil
	}
	return nil, errors.New("address not found")
}
