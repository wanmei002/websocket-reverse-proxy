package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wanmei002/websocket-reverse-proxy/tcp_server"
	"io"
	"log"
	"sync"
)

func Run() error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", tcp_server.Port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Printf("dial tcp failed: %v\n", err)
		return err
	}
	ID := uuid.New().String()

	firstWrite := &tcp_server.Device{ID: ID, Type: tcp_server.DeviceTypeServer}
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
	buffer := bytes.NewBuffer([]byte{})
Loop:
	for {
		buf := make([]byte, 1024)
		l, err := conn.Read(buf)
		if err != nil {
			log.Printf("read first write failed: %v\n", err)
			return err
		}
		for i, v := range buf {
			if v == tcp_server.MessageEnd {
				buffer.Write(buf[:i])
				break Loop
			}
		}
		buffer.Write(buf[:l])
	}
	recvData := &tcp_server.OK{}
	err = json.Unmarshal(buffer.Bytes(), recvData)
	if err != nil {
		log.Printf("unmarshal first write failed: %v\n", err)
		return err
	}
	if recvData.Code < 0 {
		log.Printf("first write failed: %v\n", recvData.Code)
		return errors.New("first write failed")
	}

	log.Printf("connect success, ID: %s\n", ID)

	// recv data read to unix
	unixConn, err := DialUnix()
	if err != nil {
		log.Printf("dial unix failed: %v\n", err)
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
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
	wg.Wait()
	return nil

}
