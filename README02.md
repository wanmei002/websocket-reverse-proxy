## grpc service
> 假设 grpc service 的服务名是 Bar

1. grpc服务在通过unix域监听请求
    ```go
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
    ```
2. 建立与代理服务的tcp连接
    ```go
   // 与代理服务器建立连接
    conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
    if err != nil {
        log.Printf("dial tcp failed: %v\n", err)
        return err
    }
    // 与代理服务器通信，告诉它 是提供服务者
    firstWrite := &tcp_server.Device{ID: deviceID, Type: tcp_server.DeviceTypeServer}
    writeData, err := json.Marshal(firstWrite)
    if err != nil {
        log.Printf("marshal first write failed: %v\n", err)
        return err
    }
    _, err = conn.Write(append(writeData, tcp_server.MessageEnd))
    if err != nil {
        log.Printf("write first write failed: %v\n", err)
        return err
    }

    bufBytes, err := tcp_server.ReadData(conn)
    if err != nil {
        log.Printf("read first write failed: %v\n", err)
        return err
    }
    recvData := &tcp_server.OK{}
    err = json.Unmarshal(bufBytes, recvData)
    if err != nil {
        log.Printf("unmarshal first write failed: %v\n", err)
        return err
    }
    if recvData.Code < 0 {
        log.Printf("first write failed: %v\n", recvData.Code)
        return errors.New("first write failed")
    }
    ```
3. 请求转发到 unix 上
   ```go
   // 与本地的 unix 建立连接
   unixConn, err := net.Dial("unix", sockPath)
	if err != nil {
		log.Fatalf("failed to dial unix: %v", err)
	}
    go func() {
        defer wg.Done()
        _, err = io.Copy(conn, unixConn)
        if err != nil {
            log.Printf("io.Copy failed: %v\n", err)
            return
        }
        fmt.Println("copy conn end")
    }()
    go func() {
        defer wg.Done()
        _, err = io.Copy(unixConn, conn)
        if err != nil {
            log.Printf("io.Copy failed unixConn: %v\n", err)
            return
        }
        fmt.Println("copy unixConn end")
    }()
    return nil
   ```
