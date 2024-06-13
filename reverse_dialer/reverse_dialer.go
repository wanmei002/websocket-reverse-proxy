package reverse_dialer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"net"
)

func Dialer() func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, deviceID string) (net.Conn, error) {
		fmt.Println("reverse dialer deviceID: " + deviceID)
		dialer, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
		forwardRequest := tcp_server.Device{
			Type: tcp_server.DeviceTypeClient,
			ID:   "",
			To:   deviceID,
		}
		sendData, err := json.Marshal(forwardRequest)
		if err != nil {
			return nil, err
		}
		_, err = dialer.Write(append(sendData, tcp_server.MessageEnd))
		if err != nil {
			return nil, err
		}
		bufBytes, err := tcp_server.ReadData(dialer)
		if err != nil {
			return nil, err
		}
		forwardResponse := tcp_server.OK{}
		err = json.Unmarshal(bufBytes, &forwardResponse)
		if err != nil {
			return nil, err
		}
		if forwardResponse.Code < 0 {
			return nil, fmt.Errorf("forward response code not correct")
		}

		return dialer, nil
	}
}
