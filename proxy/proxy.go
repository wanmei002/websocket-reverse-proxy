package proxy

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"io"
	"log"
	"net"
	"sync"
)

func Run(deviceID string, unixDial net.Conn) error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("dial tcp failed: %v\n", err)
		return err
	}

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

	log.Printf("connect success, ID: %s\n", deviceID)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err = io.Copy(conn, unixDial)
		if err != nil {
			log.Printf("io.Copy failed: %v\n", err)
			return
		}
		fmt.Println("copy conn end")
	}()
	go func() {
		defer wg.Done()
		_, err = io.Copy(unixDial, conn)
		if err != nil {
			log.Printf("io.Copy failed unixConn: %v\n", err)
			return
		}
		fmt.Println("copy unixConn end")
	}()
	wg.Wait()
	return nil

}
