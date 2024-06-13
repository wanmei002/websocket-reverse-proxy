package tcp_server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/wanmei002/websocket-reverse-proxy/private_keys"
	"io"
	"log"
	"net"
	"sync"
)

const (
	Port                  = 21112
	MessageEnd       byte = '@'
	DeviceTypeServer      = 1
	DeviceTypeClient      = 2
)

var (
	connMap     = make(map[string]net.Conn)
	connMapLock sync.RWMutex
)

type Device struct {
	Type int    `json:"type"`
	ID   string `json:"id"`
	To   string `json:"to"`
}

type OK struct {
	Code int `json:"code"`
}

func Run() error {
	ca, err := private_keys.NewCA("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	serverCA, err := tls.X509KeyPair(ca.CertPem(), ca.KeyPem())
	if err != nil {
		log.Fatal(err)
	}

	tlsServerConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCA},
	}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", Port), tlsServerConfig)
	if err != nil {
		log.Printf("Error tcp listening on port %d: %v\n", Port, err)
		return err
	}
	fmt.Println("TCP Listening on port ", Port, "; successfully")
	wg := sync.WaitGroup{}
	wg.Add(1)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			return err
		}
		fmt.Println("tcp new connection")
		go firstCommunication(conn)
	}

	wg.Wait()
	return nil
}

func firstCommunication(conn net.Conn) error {
	bufBytes, err := ReadData(conn)
	if err != nil {
		return err
	}
	device := &Device{}
	err = json.Unmarshal(bufBytes, device)
	if err != nil {
		log.Printf("Error unmarshalling json: %v\n", err)
		return err
	}
	switch device.Type {
	case DeviceTypeServer:
		err = server(device.ID, conn)
	case DeviceTypeClient:
		err = client(device.To, conn)
	default:
		err = fmt.Errorf("unknown device type: %d", device.Type)
	}
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	return err
}

func client(deviceID string, conn net.Conn) error {
	toConn := GetConn(deviceID)
	successCode := 0
	if toConn == nil {
		successCode = -1
	}
	sendData, _ := json.Marshal(OK{Code: successCode})
	_, err := conn.Write(append(sendData, MessageEnd))
	if err != nil {
		log.Printf("Error writing to connection: %v\n", err)
	}
	if successCode < 0 {
		conn.Close()
		return err
	}

	go func() {
		_, err := io.Copy(toConn, conn)
		if err != nil {
			log.Printf("client toConn error reading from connection: %v\n", err)
		}
		log.Printf("toConn conn closed.\n")
		return
	}()

	go func() {
		_, err := io.Copy(conn, toConn)
		if err != nil {
			log.Printf("client conn error reading from connection: %v\n", err)
		}
		log.Printf("conn toConn closed.\n")
	}()
	fmt.Println("Client connected to device", deviceID)
	return nil
}

func server(deviceID string, conn net.Conn) error {
	sendData, err := json.Marshal(OK{Code: 1})
	if err != nil {
		log.Printf("Error marshalling json: %v\n", err)
		return err
	}
	sendData = append(sendData, '@')
	_, err = conn.Write(sendData)
	if err != nil {
		log.Printf("Error writing to connection: %v\n", err)
		return err
	}
	setConnMap(deviceID, conn)
	return nil
}

func setConnMap(id string, conn net.Conn) {
	connMapLock.Lock()
	defer connMapLock.Unlock()

	connMap[id] = conn
}

func GetConn(id string) net.Conn {
	connMapLock.RLock()
	defer connMapLock.RUnlock()

	return connMap[id]
}

func ReadData(conn net.Conn) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
Loop:
	for {
		buf := make([]byte, 1024)
		l, err := conn.Read(buf)
		if err != nil {
			log.Printf("read first write failed: %v\n", err)
			return nil, err
		}
		for i, v := range buf {
			if v == MessageEnd {
				buffer.Write(buf[:i])
				break Loop
			}
		}
		buffer.Write(buf[:l])
	}

	return buffer.Bytes(), nil
}
