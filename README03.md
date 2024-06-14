### 写 grpc WithContextDialer 
1. 与代理服务建立连接
2. 告诉代理与哪个服务通信
```go
func Dialer() func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, deviceID string) (net.Conn, error) {
		fmt.Println("reverse dialer deviceID: " + deviceID)
		// 与代理建立连接
		dialer, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
		// 告诉代理要跟哪个服务通信
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
        // 建立连接成功
		return dialer, nil
	}
}
```

3. 使用 WithContextDialer
    ```go
    func New() (*Service, error) {
   // 使用的grpc包的版本是google.golang.org/grpc v1.64.0，最新拨号连接改成用 NewClient
    conn, err := grpc.NewClient(
   // 协议最好使用passthrough，要不然默认的使用的是 unix
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
    ```
   
4. 这样就可以通信了
```go
func (svc *Service) GetServerAddress(ctx context.Context, in *emptypb.Empty) (*messages.GetAddressResponse, error) {
	return svc.barClient.GetAddress(ctx, in)
}
```
结果返回
```shell
curl -X GET http://127.0.0.1:21114/foo/api/v1/get_address
{"address":"172.17.0.1"}
```
